package handlers

import (
	"context"
	"errors"
	"github.com/Killazius/L0/internal/domain"
	"github.com/Killazius/L0/internal/lib/api/response"
	"github.com/Killazius/L0/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
)

type OrderService interface {
	GetOrder(ctx context.Context, uid string) (*domain.Order, error)
}

type Handler struct {
	log     *zap.SugaredLogger
	service OrderService
}

func New(log *zap.SugaredLogger, service OrderService) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

// GetOrder godoc
// @Summary Get order by UID
// @Description Get order details by order UID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order_uid path string true "Order UID"
// @Success 200 {object} domain.Order "Order details"
// @Failure 400 {object} response.ErrorResponse "Invalid order UID"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /order/{order_uid} [get]
func (h *Handler) GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUID := chi.URLParam(r, "order_uid")
		if orderUID == "" {
			h.log.Info("alias is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.NewErrorResponse("order UID is required", http.StatusBadRequest, "Order UID parameter is missing"))
			return
		}
		log := h.log.With("order_uid", orderUID)

		order, err := h.service.GetOrder(r.Context(), orderUID)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrOrderNotFound):
				log.Infow("order not found")
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.NewErrorResponse("order not found", http.StatusNotFound, "The requested order was not found in the system"))
			case errors.Is(err, service.ErrInvalidOrderData):
				log.Warnw("invalid order data", "error", err)
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.NewErrorResponse("invalid order data", http.StatusBadRequest, "The requested order was not found in the system"))
			default:
				log.Errorw("internal server error", "error", err)
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.NewErrorResponse("internal server error", http.StatusInternalServerError, "The requested order was not found in the system"))
			}
			return
		}
		if order == nil {
			log.Infow("order is nil")
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.NewErrorResponse("order not found", http.StatusNotFound, "The requested order was not found in the system"))
			return
		}
		log.Info("success get order")
		render.JSON(w, r, order)
	}
}
