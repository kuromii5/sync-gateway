package routes

import (
	"github.com/gorilla/mux"
	authhandlers "github.com/kuromii5/sync-gateway/internal/handlers/auth"
	"google.golang.org/grpc"
)

// Auth service
func RegisterAuthRoutes(r *mux.Router, grpcEndpoint string, opts []grpc.DialOption) {
	r.HandleFunc("/oauth/callback", authhandlers.ExchangeCodeForTokenHandler(grpcEndpoint, opts)).Methods("GET")
}
