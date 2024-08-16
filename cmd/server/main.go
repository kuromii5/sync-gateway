package main

import "github.com/kuromii5/sync-gateway/internal/server"

func main() {
	server := server.NewServer()
	server.Run()
}
