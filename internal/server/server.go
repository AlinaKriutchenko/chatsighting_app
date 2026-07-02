package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlinaKriutchenko/chatsighting/internal/config"
	"github.com/AlinaKriutchenko/chatsighting/internal/twitch"
)

type Server struct {
	config *config.Config
	twitch *twitch.Client
}

func New(cfg *config.Config) (*Server, error) {
	tc, err := twitch.NewClient(cfg.TwitchClientID, cfg.TwitchClientSecret)
	if err != nil {
		return nil, err
	}
	return &Server{config: cfg, twitch: tc}, nil
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/vods", s.corsMiddleware(s.handleVODs))
	mux.HandleFunc("/api/analysis", s.corsMiddleware(s.rateLimitAnalysis(s.handleAnalysis)))
	mux.Handle("/", http.FileServer(http.Dir("frontend")))
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, mux)
}
