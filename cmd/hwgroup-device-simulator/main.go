package main

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log/level"
	pconfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/promlog"
	promlogflag "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/device"
)

func run(o, e io.Writer, argv ...string) int {
	app := kingpin.New(filepath.Base(argv[0]), "HWg-Push device simulator")
	promlogConfig := &promlog.Config{}

	endpoint := app.Flag("push.endpoint", "URL to send the report to.").
		Default("http://localhost:9436/portal.php").URL()
	username := app.Flag("auth.username", "Authentication principal.").
		String()
	password := app.Flag("auth.password", "Authentication secret.").
		String()

	promlogflag.AddFlags(app, promlogConfig)
	app.Version(version.Print("hwg_push device simulator"))
	app.HelpFlag.Short('h')

	if _, err := app.Parse(argv[1:]); err != nil {
		app.Errorf("%s, try --help", err)
		return 1
	}

	logger := promlog.New(promlogConfig)

	basicAuth := &pconfig.BasicAuth{
		Username: *username,
		Password: pconfig.Secret(*password),
	}
	connector := &pconfig.HTTPClientConfig{
		BasicAuth: basicAuth,
	}

	factory := device.NewNowFactory(*endpoint)
	pusher := device.NewPusher(logger, connector, factory)

	if err := pusher.Run(time.Second); err != nil {
		level.Error(logger).Log("msg", "Failed to run the reporter", "err", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Args...))
}
