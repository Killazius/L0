package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/service"
	"github.com/Killazius/L0/internal/transport/rest/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	Server *http.Server
	log    *zap.SugaredLogger
}

func NewServer(
	log *zap.SugaredLogger,
	service service.OrderService,
	cfg config.HTTPConfig,
) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)
	r.Use(newLogger(log))
	r.Use(middleware.Recoverer)

	handler := handlers.New(log, service)
	r.Route("/order", func(r chi.Router) {
		r.Get("/{order_uid}", handler.GetOrder())
	})
	swaggerURL := fmt.Sprintf("http://localhost:%s/swagger/doc.json", cfg.Port)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerURL)))

	return &Server{
		log: log,
		Server: &http.Server{
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
	s.log.Infow("rest server started", "addr", s.Server.Addr)
	err := s.Server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
func (s *Server) Close(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
