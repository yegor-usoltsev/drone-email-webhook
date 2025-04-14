package main

import (
	"context"
	"errors"
	"log"
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
		log.Printf("[INFO] server: received signal: %v, initiating shutdown", sig)
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
			log.Printf("[ERROR] server: unable to listen on %s: %v", s.addr, err)
			s.cancelServerCtx()
			return
		}
		log.Printf("[INFO] server: listening on %v", listener.Addr())
		if err = srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[ERROR] server: unable to serve: %v", err)
			s.cancelServerCtx()
		}
	}()

	<-s.serverCtx.Done()
	log.Printf("[INFO] server: shutting down")
	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), ShutdownTimeout)
	if err := srv.Shutdown(shutdownCtx); err != nil {
		cancelShutdownCtx()
		log.Fatalf("[FATAL] server: failed to shut down properly: %v", err)
	}
	cancelShutdownCtx()
	log.Printf("[INFO] server: shut down properly")
}
