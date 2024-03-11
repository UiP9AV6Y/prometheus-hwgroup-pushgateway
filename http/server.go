package http

import (
	"context"
	net "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/exporter-toolkit/web"
)

type Server struct {
	logger log.Logger
	mux    *net.ServeMux
}

func NewServer(mux *net.ServeMux, logger log.Logger) *Server {
	result := &Server{
		mux:    mux,
		logger: logger,
	}

	return result
}

func (s *Server) Serve(c *web.FlagConfig, quit <-chan struct{}) error {
	srv := &net.Server{Handler: s.mux}

	go s.wait(srv, quit)

	err := web.ListenAndServe(srv, c, s.logger)
	if err == net.ErrServerClosed {
		level.Info(s.logger).Log("msg", "HTTP server stopped")
		return nil
	}

	level.Error(s.logger).Log("msg", "HTTP server stopped", "err", err)
	return err
}

func (s *Server) wait(srv *net.Server, quit <-chan struct{}) {
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		level.Info(s.logger).Log("msg", "received termination signal")
		break
	case <-quit:
		level.Warn(s.logger).Log("msg", "received termination request")
		break
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		level.Warn(s.logger).Log("msg", "unable to shutdown cleanly", "err", err)
	}
}
