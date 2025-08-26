package handlers

import (
	"errors"
	"github.com/Killazius/L0/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	log     *zap.SugaredLogger
	service service.OrderService
}

func New(log *zap.SugaredLogger, service service.OrderService) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUID := chi.URLParam(r, "order_uid")
		if orderUID == "" {
			h.log.Info("alias is empty")
			render.Status(r, http.StatusBadRequest)
			return
		}
		log := h.log.With("order_uid", orderUID)

		order, err := h.service.GetOrder(r.Context(), orderUID)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrOrderNotFound):
				log.Infow("order not found")
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "order not found"})
			case errors.Is(err, service.ErrInvalidOrderData):
				log.Warnw("invalid order data", "error", err)
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, map[string]string{"error": "invalid order data"})
			default:
				log.Errorw("internal server error", "error", err)
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, map[string]string{"error": "internal server error"})
			}
			return
		}
		if order == nil {
			log.Infow("order is nil")
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "order not found"})
			return
		}
		log.Infow("got order", "order", order)
		render.JSON(w, r, order)
	}
}
