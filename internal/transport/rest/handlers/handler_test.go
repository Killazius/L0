package handlers

import (
	"context"
	"errors"
	"github.com/Killazius/L0/internal/domain"
	"github.com/Killazius/L0/internal/lib/api/response"
	"github.com/Killazius/L0/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetOrder(t *testing.T) {
	t.Parallel()

	testOrder := &domain.Order{
		OrderUID:    "test-uid",
		TrackNumber: "test-track",
		Entry:       "WBIL",
	}

	tests := []struct {
		name           string
		orderUID       string
		setupMock      func(*MockOrderService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "success",
			orderUID: "test-uid",
			setupMock: func(mockService *MockOrderService) {
				mockService.On("GetOrder", mock.Anything, "test-uid").
					Return(testOrder, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   testOrder,
		},
		{
			name:     "empty order uid",
			orderUID: "",
			setupMock: func(_ *MockOrderService) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   response.NewErrorResponse("order UID is required", http.StatusBadRequest, "Order UID parameter is missing"),
		},
		{
			name:     "order not found",
			orderUID: "not-found",
			setupMock: func(mockService *MockOrderService) {
				mockService.On("GetOrder", mock.Anything, "not-found").
					Return(nil, service.ErrOrderNotFound).
					Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   response.NewErrorResponse("order not found", http.StatusNotFound, "The requested order was not found in the system"),
		},
		{
			name:     "invalid order data",
			orderUID: "invalid-data",
			setupMock: func(mockService *MockOrderService) {
				mockService.On("GetOrder", mock.Anything, "invalid-data").
					Return(nil, service.ErrInvalidOrderData).
					Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   response.NewErrorResponse("invalid order data", http.StatusBadRequest, "The requested order was not found in the system"),
		},
		{
			name:     "internal server error",
			orderUID: "error",
			setupMock: func(mockService *MockOrderService) {
				mockService.On("GetOrder", mock.Anything, "error").
					Return(nil, errors.New("database error")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   response.NewErrorResponse("internal server error", http.StatusInternalServerError, "The requested order was not found in the system"),
		},
		{
			name:     "service returns nil order without error",
			orderUID: "nil-order",
			setupMock: func(mockService *MockOrderService) {
				mockService.On("GetOrder", mock.Anything, "nil-order").
					Return(nil, nil).
					Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   response.NewErrorResponse("order not found", http.StatusNotFound, "The requested order was not found in the system"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := NewMockOrderService(t)
			tt.setupMock(mockService)

			logger := zap.NewNop().Sugar()
			handler := New(logger, mockService)

			req, err := http.NewRequest("GET", "/order/"+tt.orderUID, nil)
			require.NoError(t, err)

			if tt.orderUID != "" {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("order_uid", tt.orderUID)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			}

			rr := httptest.NewRecorder()

			handlerFunc := handler.GetOrder()
			handlerFunc(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, rr.Body.String(), `"order_uid":"test-uid"`)
			} else {
				assert.Contains(t, rr.Body.String(), `"error":`)
			}

			mockService.AssertExpectations(t)
		})
	}
}
func TestHandler_GetOrder_ChiContext(t *testing.T) {
	t.Parallel()

	mockService := NewMockOrderService(t)
	mockService.On("GetOrder", mock.Anything, "test-uid").
		Return(&domain.Order{OrderUID: "test-uid"}, nil).
		Once()

	logger := zap.NewNop().Sugar()
	handler := New(logger, mockService)

	req, err := http.NewRequest("GET", "/order/test-uid", nil)
	require.NoError(t, err)

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("order_uid", "test-uid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	rr := httptest.NewRecorder()
	handlerFunc := handler.GetOrder()
	handlerFunc(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetOrder_WithoutChiContext(t *testing.T) {
	t.Parallel()

	mockService := NewMockOrderService(t)

	logger := zap.NewNop().Sugar()
	handler := New(logger, mockService)

	req, err := http.NewRequest("GET", "/order/test-uid", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handlerFunc := handler.GetOrder()
	handlerFunc(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_GetOrder_ContextCancellation(t *testing.T) {
	t.Parallel()

	mockService := NewMockOrderService(t)
	mockService.On("GetOrder", mock.Anything, "test-uid").
		Return(nil, context.Canceled).
		Once()

	logger := zap.NewNop().Sugar()
	handler := New(logger, mockService)

	req, err := http.NewRequest("GET", "/order/test-uid", nil)
	require.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("order_uid", "test-uid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handlerFunc := handler.GetOrder()
	handlerFunc(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "internal server error")
	mockService.AssertExpectations(t)
}
