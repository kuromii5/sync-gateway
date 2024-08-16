package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kuromii5/sync-gateway/internal/config"
	"github.com/kuromii5/sync-gateway/internal/server/gateway"
	"github.com/kuromii5/sync-gateway/internal/server/logger"
)

type Server struct {
	gateway *gateway.Gateway
}

func NewServer() *Server {
	config := config.Load()

	logger := logger.New(config.Env, config.LogLevel)

	endpoints := map[string]string{
		"auth":         config.AuthServiceEndpoint,
		"user":         config.UserServiceEndpoint,
		"feed":         config.FeedServiceEndpoint,
		"message":      config.MessageServiceEndpoint,
		"notification": config.NotificationServiceEndpoint,
		"music":        config.MusicServiceEndpoint,
		"video":        config.VideoServiceEndpoint,
		"group":        config.GroupServiceEndpoint,
	}

	gateway := gateway.NewGateway(config.Port, logger, endpoints)

	return &Server{gateway: gateway}
}

func (s *Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		s.gateway.Run(ctx)
	}()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer shutdownCancel()

	cancel()
	s.shutdown(shutdownCtx)
}

func (s *Server) shutdown(ctx context.Context) {
	s.gateway.Shutdown(ctx)
}
