package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	vercol "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/route"
	"github.com/prometheus/common/version"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"

	promlogflag "github.com/prometheus/common/promlog/flag"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/collector"
	www "github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/http"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/http/handler"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/log/gokitlog"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao"
	portalflag "github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/flag"
)

func run(argv []string) (exitCode int) {
	var (
		app           = kingpin.New(filepath.Base(argv[0]), "HWg-Push protocol gateway")
		promlogConfig = &promlog.Config{}
		portalConfig  = &portal.Config{}
		webConfig     = webflag.AddFlags(app, ":9436")

		disableGatewayMetrics = app.Flag("web.disable-gateway-metrics",
			"Exclude metrics about the gateway itself (promhttp_*, process_*, go_*).").
			Bool()
		maxRequests = app.Flag("web.max-requests",
			"Maximum number of parallel scrape requests. Use 0 to disable.").
			Default("40").Int()
		routePrefix = app.Flag("web.route-prefix",
			"Prefix for the internal routes of web endpoints.").
			Default("").String()
		metricsPath = app.Flag("web.telemetry-path",
			"Path under which to expose metrics.").
			Default("/metrics").String()

		persistenceFile = app.Flag("persistence.file",
			"Filesystem location to persist device states").
			Default("").String()
		persistenceInterval = app.Flag("persistence.interval",
			"The minimum interval at which to write out the persistence file.").
			Default("5m").Duration()

		portalPath = app.Flag("portal.report-path",
			"Path under which to receive push protocol reports.").
			Default("/portal.php").String()
		dumpPath = app.Flag("portal.dump-path",
			"Path under which to query a dump of the internal data state.").
			Default("").String()
	)

	promlogflag.AddFlags(app, promlogConfig)
	portalflag.AddFlags(app, portalConfig)
	app.Version(version.Print(collector.Namespace))
	app.HelpFlag.Short('h')

	exitCode = 1 // assume error until the very end

	if _, err := app.Parse(argv[1:]); err != nil {
		app.Errorf("%s, try --help", err)
		return
	}

	logger := promlog.New(promlogConfig)
	persistence := dao.New(logger)

	if db := *persistenceFile; db != "" {
		dbSave := func() {
			level.Info(logger).Log("msg", "saving database", "file", db)

			if err := persistence.ExportFile(db); err != nil {
				level.Error(logger).Log("msg", "unable to save database", "err", err)
				exitCode = 1
			}
		}
		dbTick := time.Tick(*persistenceInterval)
		dbSaverTick := func() {
			for _ = range dbTick {
				dbSave()
			}
		}

		if err := persistence.ImportFile(db); err != nil {
			level.Error(logger).Log("msg", "unable to restore database", "err", err)
			return
		}

		defer dbSave()

		if dbTick != nil {
			go dbSaverTick()
		}
	}

	processor, err := portal.New(portalConfig, persistence, logger)
	if err != nil {
		level.Error(logger).Log("msg", "unable to create push protocol processor", "err", err)
		return
	}

	gatewayMetricsRegistry := prometheus.NewRegistry()
	if err := gatewayMetricsRegistry.Register(vercol.NewCollector(collector.Namespace)); err != nil {
		level.Error(logger).Log("msg", "unable to register version metrics", "err", err)
		return
	}

	if err := handler.RegisterMetrics(gatewayMetricsRegistry); err != nil {
		level.Error(logger).Log("msg", "unable to register handler metrics", "err", err)
		return
	}

	deviceMetricsRegistry := prometheus.NewRegistry()
	if err := deviceMetricsRegistry.Register(collector.New(persistence, logger)); err != nil {
		level.Error(logger).Log("msg", "unable to register collector", "err", err)
		return
	}

	gatherers := prometheus.Gatherers{
		deviceMetricsRegistry,
	}
	if !*disableGatewayMetrics {
		gatherers = append(gatherers,
			prometheus.DefaultGatherer,
			gatewayMetricsRegistry,
		)
	}

	logProxy := gokitlog.NewLoggerProxy(level.Error(logger), "msg")
	handlerOpts := promhttp.HandlerOpts{
		ErrorLog:            logProxy,
		ErrorHandling:       promhttp.ContinueOnError,
		MaxRequestsInFlight: *maxRequests,
	}
	metricsHandler := promhttp.HandlerFor(gatherers, handlerOpts)
	dumpHandler := handler.Dump(processor, logger)
	pushHandler := handler.Push(processor, logger)
	router := route.New().WithPrefix(*routePrefix)
	quitCh := make(chan struct{})

	if !*disableGatewayMetrics {
		metricsHandler = promhttp.InstrumentMetricHandler(
			gatewayMetricsRegistry, metricsHandler,
		)
		dumpHandler = handler.InstrumentWithSummaries(
			handler.InstrumentWithCounter("dump", dumpHandler),
		)
		pushHandler = handler.InstrumentWithSummaries(
			handler.InstrumentWithCounter("push", pushHandler),
		)
	}

	router.Post(*portalPath, pushHandler)
	router.Get(*metricsPath, metricsHandler.ServeHTTP)

	if dump := *dumpPath; dump != "" {
		router.Get(dump, dumpHandler)
	}

	mux := http.NewServeMux()
	mux.Handle("/", router)

	level.Info(logger).Log("msg", "starting hwg_gateway", "version", version.Info())
	level.Info(logger).Log("build_context", version.BuildContext())

	server := www.NewServer(mux, logger)
	if err := server.Serve(webConfig, quitCh); err != nil {
		return
	}

	exitCode = 0

	return
}

func main() {
	os.Exit(run(os.Args))
}
