package app

import (
	"context"
	"log"
	"yema.dev/pkg/server"
)

type App struct {
	name    string
	servers []server.Server
}

type Option func(a *App)

func NewApp(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithServer(servers ...server.Server) Option {
	return func(a *App) {
		a.servers = servers
	}
}

func WithName(name string) Option {
	return func(a *App) {
		a.name = name
	}
}

func (a *App) Run(ctx context.Context) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	for _, srv := range a.servers {
		go func(srv server.Server) {
			err := srv.Start(ctx)
			if err != nil {
				log.Printf("Server start err: %v", err)
			}
		}(srv)
	}

	select {
	case <-ctx.Done():
		// Context canceled
		log.Println("Context canceled")
	}

	// Gracefully stop the servers
	for _, srv := range a.servers {
		err := srv.Stop(ctx)
		if err != nil {
			log.Printf("Server stop err: %v", err)
		}
	}

	return nil
}
