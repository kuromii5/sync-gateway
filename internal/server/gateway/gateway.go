package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	auth "github.com/kuromii5/sync-auth/api/sync-auth/v1"
	"github.com/kuromii5/sync-gateway/internal/routes"
	le "github.com/kuromii5/sync-gateway/pkg/logger/l_err"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type Gateway struct {
	port   int
	logger *slog.Logger

	authEndpoint         string
	userEndpoint         string
	feedEndpoint         string
	messageEndpoint      string
	notificationEndpoint string
	musicEndpoint        string
	videoEndpoint        string
	groupEndpoint        string
}

func NewGateway(port int, logger *slog.Logger, endpoints map[string]string) *Gateway {
	return &Gateway{
		port:                 port,
		logger:               logger,
		authEndpoint:         endpoints["auth"],
		userEndpoint:         endpoints["user"],
		feedEndpoint:         endpoints["feed"],
		messageEndpoint:      endpoints["message"],
		notificationEndpoint: endpoints["notification"],
		musicEndpoint:        endpoints["music"],
		videoEndpoint:        endpoints["video"],
		groupEndpoint:        endpoints["group"],
	}
}

func (g *Gateway) Run(ctx context.Context) {
	r := mux.NewRouter()

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// auth service
	err := auth.RegisterAuthHandlerFromEndpoint(ctx, mux, g.authEndpoint, opts)
	if err != nil {
		g.logger.Error("Failed to register Auth service gateway", le.Err(err))
	}

	routes.RegisterAuthRoutes(r, g.authEndpoint, opts)

	r.PathPrefix("/").Handler(mux)

	handlerWithCORS := corsMiddleware(r)

	addr := fmt.Sprintf(":%d", g.port)
	g.logger.Info("Starting API gateway...", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, handlerWithCORS); err != nil {
		g.logger.Error("Failed to start API gateway", le.Err(err))
	}
}

func (g *Gateway) Shutdown(ctx context.Context) {
	g.logger.Info("Shutting down API gateway...")

	<-ctx.Done()

	g.logger.Info("API gateway successfully shut down")
}
