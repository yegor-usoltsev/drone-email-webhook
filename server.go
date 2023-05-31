package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	ReadTimeout       = 16 * time.Second
	ReadHeaderTimeout = 5 * time.Second
	HandlerTimeout    = 10 * time.Second
	WriteTimeout      = 11 * time.Second
	IdleTimeout       = 120 * time.Second
	ShutdownTimeout   = 11 * time.Second
	MaxHeaderBytes    = 16 * 1024   // 16 KB
	MaxBodyBytes      = 1024 * 1024 // 1 MB
)

type Server struct {
	addr            string
	serverCtx       context.Context
	cancelServerCtx func()
}

func NewServer(settings Settings) *Server {
	addr := fmt.Sprintf("%s:%d", settings.ServerHost, settings.ServerPort)
	serverCtx, cancelServerCtx := context.WithCancel(context.Background())
	return &Server{addr: addr, serverCtx: serverCtx, cancelServerCtx: cancelServerCtx}
}

func (s *Server) ListenAndServe(handler http.Handler) {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, os.Kill)
	go func() {
		log.Println("Server received signal:", <-shutdownCh)
		s.cancelServerCtx()
	}()

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
			log.Println("Unable to listen:", err)
			s.cancelServerCtx()
			return
		}
		log.Println("Server is listening on", listener.Addr())
		if err = srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Println("Unable to serve:", err)
			s.cancelServerCtx()
		}
	}()

	<-s.serverCtx.Done()
	log.Println("Server is shutting down")
	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancelShutdownCtx()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalln("Server failed to shut down properly:", err)
	}
	log.Println("Server shut down properly")
}
