package transport

import (
	"github.com/Killazius/L0/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
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
		order, err := h.service.GetOrder(r.Context(), orderUID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, order)
	}
}
