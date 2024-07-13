package server

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/AxelTahmid/golang-starter/internal/middlewares"
	"github.com/AxelTahmid/golang-starter/internal/modules/auth"
)

func (s *Server) routes() {
	// global middlewares
	s.router.Use(chiMiddleware.RealIP)
	s.router.Use(chiMiddleware.RequestID)
	s.router.Use(middlewares.Logger(s.log))
	s.router.Use(middlewares.Recovery)
	s.router.Use(middlewares.Secure(s.conf.Secure).Handler)
	s.router.Use(middlewares.Json)
	s.router.Use(chiMiddleware.Heartbeat("/ping"))

	// routes
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/auth", auth.Routes(s.db))
	})
}
