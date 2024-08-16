package auth

import (
	"context"
	"encoding/json"
	"net/http"

	auth "github.com/kuromii5/sync-auth/api/sync-auth/v1"
	"google.golang.org/grpc"
)

func ExchangeCodeForTokenHandler(grpcEndpoint string, opts []grpc.DialOption) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		provider := r.URL.Query().Get("provider")
		if provider == "" {
			http.Error(w, "Provider not found", http.StatusBadRequest)
			return
		}

		conn, err := grpc.NewClient(grpcEndpoint, opts...)
		if err != nil {
			http.Error(w, "Failed to connect to gRPC server: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := auth.NewAuthClient(conn)
		response, err := client.ExchangeCodeForToken(context.Background(), &auth.ExchangeCodeRequest{
			Provider: provider,
			Code:     code,
		})
		if err != nil {
			http.Error(w, "Failed to exchange code: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
