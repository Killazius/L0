package service

import (
	"context"
	"errors"
	"github.com/Killazius/L0/internal/lib/test"
	"github.com/Killazius/L0/internal/repository"
	"testing"
	"time"

	"github.com/Killazius/L0/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestService_GetOrder(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)
	defer zap.ReplaceGlobals(zap.NewNop())

	testOrder := test.GenerateOrder()

	tests := []struct {
		name          string
		setupMocks    func(*MockOrderRepository, *MockOrderCache)
		orderUID      string
		expectedOrder *domain.Order
		expectedError error
	}{
		{
			name: "success from cache",
			setupMocks: func(_ *MockOrderRepository, cache *MockOrderCache) {
				cache.On("Get", mock.Anything, "test-uid").
					Return(testOrder, nil).
					Once()
			},
			orderUID:      "test-uid",
			expectedOrder: testOrder,
			expectedError: nil,
		},
		{
			name: "success from database with cache miss",
			setupMocks: func(repo *MockOrderRepository, cache *MockOrderCache) {
				cache.On("Get", mock.Anything, "test-uid").
					Return(nil, repository.ErrOrderNotFound).
					Once()
				repo.On("Get", mock.Anything, "test-uid").
					Return(testOrder, nil).
					Once()
				cache.On("Set", mock.Anything, testOrder).
					Return(nil).
					Once()
			},
			orderUID:      "test-uid",
			expectedOrder: testOrder,
			expectedError: nil,
		},
		{
			name: "cache error falls back to database success",
			setupMocks: func(repo *MockOrderRepository, cache *MockOrderCache) {
				cache.On("Get", mock.Anything, "test-uid").
					Return(nil, errors.New("cache error")).
					Once()
				repo.On("Get", mock.Anything, "test-uid").
					Return(testOrder, nil).
					Once()
				cache.On("Set", mock.Anything, testOrder).
					Return(nil).
					Once()
			},
			orderUID:      "test-uid",
			expectedOrder: testOrder,
			expectedError: nil,
		},
		{
			name: "order not found",
			setupMocks: func(repo *MockOrderRepository, cache *MockOrderCache) {
				cache.On("Get", mock.Anything, "not-found").
					Return(nil, repository.ErrOrderNotFound).
					Once()
				repo.On("Get", mock.Anything, "not-found").
					Return(nil, repository.ErrOrderNotFound).
					Once()
			},
			orderUID:      "not-found",
			expectedOrder: nil,
			expectedError: ErrOrderNotFound,
		},
		{
			name: "database returns invalid data error",
			setupMocks: func(repo *MockOrderRepository, cache *MockOrderCache) {
				cache.On("Get", mock.Anything, "invalid-data").
					Return(nil, repository.ErrOrderNotFound).
					Once()
				repo.On("Get", mock.Anything, "invalid-data").
					Return(nil, repository.ErrDeliveryNotFound).
					Once()
			},
			orderUID:      "invalid-data",
			expectedOrder: nil,
			expectedError: ErrInvalidOrderData,
		},
		{
			name: "database returns unexpected error",
			setupMocks: func(repo *MockOrderRepository, cache *MockOrderCache) {
				cache.On("Get", mock.Anything, "db-error").
					Return(nil, repository.ErrOrderNotFound).
					Once()
				repo.On("Get", mock.Anything, "db-error").
					Return(nil, errors.New("database connection failed")).
					Once()
			},
			orderUID:      "db-error",
			expectedOrder: nil,
			expectedError: errors.New("failed to get order from database"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockOrderRepository(t)
			mockCache := NewMockOrderCache(t)
			tt.setupMocks(mockRepo, mockCache)

			service := New(mockRepo, mockCache)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := service.GetOrder(ctx, tt.orderUID)
			service.wg.Wait()
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOrder, result)
			}

			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestService_CreateOrder(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)
	defer zap.ReplaceGlobals(zap.NewNop())

	validOrder := test.GenerateOrder()

	invalidOrder := &domain.Order{
		OrderUID: "",
	}

	tests := []struct {
		name          string
		setupMocks    func(*MockOrderRepository, *MockOrderCache)
		order         *domain.Order
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func(repo *MockOrderRepository, cache *MockOrderCache) {
				repo.On("Create", mock.Anything, validOrder).
					Return(nil).
					Once()
				cache.On("Set", mock.Anything, validOrder).
					Return(nil).
					Once()
			},
			order:         validOrder,
			expectedError: nil,
		},
		{
			name: "invalid order data",
			setupMocks: func(_ *MockOrderRepository, _ *MockOrderCache) {
			},
			order:         invalidOrder,
			expectedError: ErrInvalidOrderData,
		},
		{
			name: "nil order",
			setupMocks: func(_ *MockOrderRepository, _ *MockOrderCache) {
			},
			order:         nil,
			expectedError: ErrInvalidOrderData,
		},
		{
			name: "order already exists",
			setupMocks: func(repo *MockOrderRepository, _ *MockOrderCache) {
				repo.On("Create", mock.Anything, validOrder).
					Return(repository.ErrDuplicateOrder).
					Once()
			},
			order:         validOrder,
			expectedError: ErrOrderAlreadyExists,
		},
		{
			name: "database error",
			setupMocks: func(repo *MockOrderRepository, _ *MockOrderCache) {
				repo.On("Create", mock.Anything, validOrder).
					Return(errors.New("database error")).
					Once()
			},
			order:         validOrder,
			expectedError: errors.New("failed to create order"),
		},
		{
			name: "cache set error after successful create",
			setupMocks: func(repo *MockOrderRepository, cache *MockOrderCache) {
				repo.On("Create", mock.Anything, validOrder).
					Return(nil).
					Once()
				cache.On("Set", mock.Anything, validOrder).
					Return(errors.New("cache error")).
					Once()
			},
			order:         validOrder,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockOrderRepository(t)
			mockCache := NewMockOrderCache(t)
			tt.setupMocks(mockRepo, mockCache)

			service := New(mockRepo, mockCache)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := service.CreateOrder(ctx, tt.order)
			service.wg.Wait()

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestService_ContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("GetOrder context cancellation", func(t *testing.T) {
		t.Parallel()

		mockRepo := NewMockOrderRepository(t)
		mockCache := NewMockOrderCache(t)

		mockCache.On("Get", mock.Anything, "test-uid").
			Return(func(ctx context.Context, _ string) (*domain.Order, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(100 * time.Millisecond):
					return nil, repository.ErrOrderNotFound
				}
			}).
			Once()

		service := New(mockRepo, mockCache)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := service.GetOrder(ctx, "test-uid")

		require.Error(t, err)
		assert.Contains(t, err.Error(), context.DeadlineExceeded.Error())
	})

	t.Run("CreateOrder context cancellation during cache set", func(t *testing.T) {
		t.Parallel()

		mockRepo := NewMockOrderRepository(t)
		mockCache := NewMockOrderCache(t)

		validOrder := test.GenerateOrder()

		mockRepo.On("Create", mock.Anything, validOrder).
			Return(nil).
			Once()

		mockCache.On("Set", mock.Anything, validOrder).
			Return(context.Canceled).
			Once()

		service := New(mockRepo, mockCache)

		ctx := context.Background()
		err := service.CreateOrder(ctx, validOrder)

		require.NoError(t, err)

		time.Sleep(50 * time.Millisecond)

		mockRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})
}

func TestService_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("nil order in CreateOrder", func(t *testing.T) {
		t.Parallel()

		service := New(NewMockOrderRepository(t), NewMockOrderCache(t))
		err := service.CreateOrder(context.Background(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), ErrInvalidOrderData.Error())
	})

	t.Run("empty order UID in GetOrder", func(t *testing.T) {
		t.Parallel()

		mockRepo := NewMockOrderRepository(t)
		mockCache := NewMockOrderCache(t)

		mockCache.On("Get", mock.Anything, "").
			Return(nil, repository.ErrOrderNotFound).
			Once()
		mockRepo.On("Get", mock.Anything, "").
			Return(nil, repository.ErrOrderNotFound).
			Once()

		service := New(mockRepo, mockCache)
		_, err := service.GetOrder(context.Background(), "")

		require.Error(t, err)
		assert.Contains(t, err.Error(), ErrOrderNotFound.Error())
	})
}
