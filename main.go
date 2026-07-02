package main

import (
	"log"

	"github.com/AlinaKriutchenko/chatsighting/internal/config"
	"github.com/AlinaKriutchenko/chatsighting/internal/server"
)

func main() {
	cfg := config.Load()

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
