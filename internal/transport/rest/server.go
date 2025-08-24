package rest

import (
	"context"
	"errors"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/service"
	"github.com/Killazius/L0/internal/transport/rest/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	server *http.Server
	log    *zap.SugaredLogger
}

func NewServer(
	log *zap.SugaredLogger,
	service *service.Service,
	cfg config.HTTPConfig,
) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(newLogger(log))

	handler := handlers.New(log, service)

	r.Route("/order", func(r chi.Router) {
		r.Get("/{order_uid}", handler.GetOrder())
	})

	return &Server{
		log: log,
		server: &http.Server{
			Addr:         cfg.GetAddr(),
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
			Handler:      r,
		},
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		s.log.Fatalw("failed to run HTTP-server", "error", err)
	}
}
func (s *Server) Run() error {
	err := s.server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
func (s *Server) Stop(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Errorw("failed to stop HTTP server", "error", err)
	}
}
