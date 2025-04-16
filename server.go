package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	ReadTimeout       = 16 * time.Second
	ReadHeaderTimeout = 5 * time.Second
	WriteTimeout      = 11 * time.Second
	IdleTimeout       = 120 * time.Second
	ShutdownTimeout   = 11 * time.Second
	MaxHeaderBytes    = 16 * 1024   // 16 KB
	MaxBodyBytes      = 1024 * 1024 // 1 MB
)

type Server struct {
	addr            string
	serverCtx       context.Context //nolint:containedctx
	cancelServerCtx func()
}

func NewServer(settings Settings) *Server {
	addr := net.JoinHostPort(settings.ServerHost, strconv.Itoa(int(settings.ServerPort)))
	serverCtx, cancelServerCtx := context.WithCancel(context.Background())
	return &Server{addr: addr, serverCtx: serverCtx, cancelServerCtx: cancelServerCtx}
}

func (s *Server) ListenAndServe(handler http.Handler) {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-shutdownCh
		slog.Info("http server received shutdown signal", "signal", sig)
		s.cancelServerCtx()
	}()

	//nolint:exhaustruct
	srv := &http.Server{
		ReadTimeout:       ReadTimeout,
		ReadHeaderTimeout: ReadHeaderTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
		MaxHeaderBytes:    MaxHeaderBytes,
		Handler:           http.MaxBytesHandler(handler, MaxBodyBytes),
	}
	go func() {
		listener, err := net.Listen("tcp", s.addr)
		if err != nil {
			slog.Error("http server unable to listen on address", "addr", s.addr, "error", err)
			s.cancelServerCtx()
			return
		}
		slog.Info("http server started listening", "addr", listener.Addr().String())
		if err = srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server unable to handle requests", "error", err)
			s.cancelServerCtx()
		}
	}()

	<-s.serverCtx.Done()
	slog.Info("http server initiating shutdown")
	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancelShutdownCtx()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}
	slog.Info("http server completed shutdown")
}
