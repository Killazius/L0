package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"l0/internal/service"
	"net/http"
)

type Handler struct {
	log     *zap.SugaredLogger
	service *service.Service
}

func newHandler(log *zap.SugaredLogger, service *service.Service) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) Info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUID := chi.URLParam(r, "order_uid")
		if orderUID == "" {
			h.log.Info("alias is empty")
			render.Status(r, http.StatusBadRequest)
			return
		}
		render.JSON(w, r, orderUID)
	}
}
