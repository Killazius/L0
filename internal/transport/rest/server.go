package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/Killazius/L0/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	server *http.Server
	log    *zap.SugaredLogger
}

type Handler interface {
	GetOrder() http.HandlerFunc
}

func NewServer(
	log *zap.SugaredLogger,
	handler Handler,
	cfg config.HTTPConfig,
) *Server {

	return &Server{
		log: log,
		server: &http.Server{
			Addr:         cfg.GetAddr(),
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
			Handler:      registerRoutes(handler, log, cfg.Port),
		},
	}
}

func registerRoutes(h Handler, log *zap.SugaredLogger, port string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)
	r.Use(logMiddleware(log))
	r.Use(middleware.Recoverer)

	r.Route("/order", func(r chi.Router) {
		r.Get("/{order_uid}", h.GetOrder())
	})
	swaggerURL := fmt.Sprintf("http://localhost:%s/swagger/doc.json", port)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerURL)))

	r.Handle("/*", http.FileServer(http.Dir("./static")))
	return r
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		s.log.Fatalw("failed to run HTTP-server", "error", err)
	}
}
func (s *Server) Run() error {
	s.log.Infow("rest server started", "addr", s.Addr())
	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
func (s *Server) Close(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
func (s *Server) Addr() string {
	return s.server.Addr
}
