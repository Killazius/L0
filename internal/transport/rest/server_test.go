package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Killazius/L0/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop().Sugar()
	mockHandler := NewMockHandler(t)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mockHandler.On("GetOrder").Return(handlerFunc).Once()

	cfg := config.HTTPConfig{
		Host:        "localhost",
		Port:        "8080",
		Timeout:     30 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	server := NewServer(logger, mockHandler, cfg)

	assert.NotNil(t, server)
	assert.Equal(t, "localhost:8080", server.Addr())

	mockHandler.AssertExpectations(t)
}

func TestServer_Addr(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop().Sugar()
	mockHandler := NewMockHandler(t)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mockHandler.On("GetOrder").Return(handlerFunc).Once()

	cfg := config.HTTPConfig{
		Host: "127.0.0.1",
		Port: "9090",
	}

	server := NewServer(logger, mockHandler, cfg)

	assert.Equal(t, "127.0.0.1:9090", server.Addr())
	mockHandler.AssertExpectations(t)
}

func TestRegisterRoutes(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop().Sugar()
	mockHandler := NewMockHandler(t)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mockHandler.On("GetOrder").Return(handlerFunc).Once()

	router := registerRoutes(mockHandler, logger, "8080")

	tests := []struct {
		name     string
		method   string
		path     string
		expected int
	}{
		{
			name:     "order route",
			method:   "GET",
			path:     "/order/test-uid",
			expected: http.StatusOK,
		},
		{
			name:     "static files route",
			method:   "GET",
			path:     "/some-file.txt",
			expected: http.StatusNotFound,
		},
		{
			name:     "not found route",
			method:   "GET",
			path:     "/non-existent",
			expected: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(tt.method, tt.path, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expected, rr.Code)
		})
	}

	mockHandler.AssertExpectations(t)
}

func TestLogMiddleware(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop().Sugar()

	middlewareFunc := logMiddleware(logger)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("test response"))
		if err != nil {
			return
		}
	})

	wrappedHandler := middlewareFunc(testHandler)

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{}))

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "test response", rr.Body.String())
}

func TestServer_Close(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop().Sugar()
	mockHandler := NewMockHandler(t)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mockHandler.On("GetOrder").Return(handlerFunc).Once()

	cfg := config.HTTPConfig{
		Host: "localhost",
		Port: "0",
	}

	server := NewServer(logger, mockHandler, cfg)

	go func() {
		err := server.Run()
		if err != nil {
			t.Logf("Server run error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := server.Close(ctx)
	require.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestServer_Run_Error(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop().Sugar()
	mockHandler := NewMockHandler(t)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mockHandler.On("GetOrder").Return(handlerFunc).Once()

	cfg := config.HTTPConfig{
		Host: "localhost",
		Port: "abc",
	}

	server := NewServer(logger, mockHandler, cfg)

	err := server.Run()
	require.Error(t, err)
	mockHandler.AssertExpectations(t)
}

func TestServer_MustRun(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop().Sugar()
	mockHandler := NewMockHandler(t)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mockHandler.On("GetOrder").Return(handlerFunc).Once()

	cfg := config.HTTPConfig{
		Host: "localhost",
		Port: "0",
	}

	server := NewServer(logger, mockHandler, cfg)

	assert.NotPanics(t, func() {
		go server.MustRun()
		time.Sleep(100 * time.Millisecond)
		server.Close(context.Background())
	})
	mockHandler.AssertExpectations(t)
}

func TestRegisterRoutes_OrderRoute(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop().Sugar()
	mockHandler := NewMockHandler(t)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("order response"))
		if err != nil {
			return
		}
	})
	mockHandler.On("GetOrder").Return(handlerFunc).Once()

	router := registerRoutes(mockHandler, logger, "8080")

	req, err := http.NewRequest("GET", "/order/12345", nil)
	require.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("order_uid", "12345")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "order response", rr.Body.String())

	mockHandler.AssertExpectations(t)
}

func TestServer_Close_WithTimeout(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop().Sugar()
	mockHandler := NewMockHandler(t)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mockHandler.On("GetOrder").Return(handlerFunc).Once()

	cfg := config.HTTPConfig{
		Host: "localhost",
		Port: "0",
	}

	server := NewServer(logger, mockHandler, cfg)

	listenErr := make(chan error, 1)
	go func() {
		listenErr <- server.Run()
	}()

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(2 * time.Millisecond)

	err := server.Close(ctx)
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		server.Close(ctx)
	})

	server.Close(context.Background())

	mockHandler.AssertExpectations(t)
}
