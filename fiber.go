package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Fiber struct {
	*fiber.App
	settings        Settings
	serverCtx       context.Context
	cancelServerCtx func()
}

func NewFiber(settings Settings) *Fiber {
	serverCtx, cancelServerCtx := context.WithCancel(context.Background())
	return &Fiber{
		App: fiber.New(fiber.Config{
			ReadBufferSize:        settings.ServerMaxHeaderBytes,
			BodyLimit:             settings.ServerMaxBodyBytes,
			ReadTimeout:           settings.ServerReadTimeout,
			WriteTimeout:          settings.ServerWriteTimeout,
			IdleTimeout:           settings.ServerIdleTimeout,
			DisableStartupMessage: true,
		}),
		settings:        settings,
		serverCtx:       serverCtx,
		cancelServerCtx: cancelServerCtx,
	}
}

func (f *Fiber) Context() context.Context {
	return f.serverCtx
}

func (f *Fiber) ListenGracefully() {
	go func() {
		quitCh := make(chan os.Signal, 1)
		signal.Notify(quitCh, os.Interrupt, os.Kill)
		log.Println("Server received signal:", <-quitCh)
		f.cancelServerCtx()
	}()

	go func() {
		addr := fmt.Sprintf("%s:%d", f.settings.ServerHost, f.settings.ServerPort)
		log.Println("Server is listening on http://" + strings.ReplaceAll(addr, "0.0.0.0", "localhost"))
		if err := f.Listen(addr); err != nil {
			log.Println("Unable to listen:", err)
			f.cancelServerCtx()
		}
	}()

	<-f.serverCtx.Done()
	if err := f.ShutdownWithTimeout(f.settings.ServerShutdownTimeout); err != nil {
		log.Panicln("Server failed to shut down properly:", err)
	}
	log.Println("Server shut down properly")
}
