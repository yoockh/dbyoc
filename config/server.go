package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yoockh/dbyoc/config"
)

// Server wraps an http.Server with graceful start/stop helpers.
type Server struct {
	httpServer *http.Server
	cfg        config.ServerConfig
	logger     *logrus.Logger
}

// New creates a new Server. Handler is the http.Handler (router).
func New(cfg config.ServerConfig, handler http.Handler, logger *logrus.Logger) *Server {
	// sensible defaults if zero
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 5 // seconds
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 10 // seconds
	}
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 10 // seconds
	}

	httpSrv := &http.Server{
		Addr:         cfg.Address(),
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.ShutdownTimeout) * time.Second,
	}

	return &Server{
		httpServer: httpSrv,
		cfg:        cfg,
		logger:     logger,
	}
}

// ListenAndServe starts the HTTP server (blocking). If TLS is configured it uses ListenAndServeTLS.
func (s *Server) ListenAndServe() error {
	if s.cfg.TLS {
		s.logger.Infof("starting TLS server on %s", s.httpServer.Addr)
		return s.httpServer.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	}
	s.logger.Infof("starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// StartWithSignals starts the server in a goroutine and listens for OS signals to gracefully shutdown.
// It blocks until shutdown completes or context is canceled.
func (s *Server) StartWithSignals(ctx context.Context) error {
	// start server
	errCh := make(chan error, 1)
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// handle signals and context cancellation
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	select {
	case <-ctx.Done():
		s.logger.Info("context canceled, shutting down server")
	case sig := <-sigCh:
		s.logger.Infof("received signal %v, shutting down server", sig)
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	}

	// perform graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.ShutdownTimeout)*time.Second)
	defer cancel()

	s.logger.Info("shutting down HTTP server")
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	s.logger.Info("server stopped gracefully")
	return nil
}

// Stop force-closes the server (immediate).
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
