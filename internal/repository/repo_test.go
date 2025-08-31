package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Killazius/L0/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestore(t *testing.T) {
	t.Parallel()

	testOrders := []domain.Order{
		{OrderUID: "order1", TrackNumber: "track1"},
		{OrderUID: "order2", TrackNumber: "track2"},
		{OrderUID: "order3", TrackNumber: "track3"},
	}

	tests := []struct {
		name          string
		setupMocks    func(*MockOrderProvider, *MockOrderSetter)
		workers       int
		expectedError string
	}{
		{
			name: "success with multiple workers",
			setupMocks: func(repo *MockOrderProvider, cache *MockOrderSetter) {
				repo.On("GetAll", mock.Anything).
					Return(testOrders, nil).
					Once()
				cache.On("Set", mock.Anything, &testOrders[0]).
					Return(nil).
					Once()
				cache.On("Set", mock.Anything, &testOrders[1]).
					Return(nil).
					Once()
				cache.On("Set", mock.Anything, &testOrders[2]).
					Return(nil).
					Once()
			},
			workers: 3,
		},
		{
			name: "success with single worker",
			setupMocks: func(repo *MockOrderProvider, cache *MockOrderSetter) {
				repo.On("GetAll", mock.Anything).
					Return(testOrders, nil).
					Once()
				cache.On("Set", mock.Anything, &testOrders[0]).
					Return(nil).
					Once()
				cache.On("Set", mock.Anything, &testOrders[1]).
					Return(nil).
					Once()
				cache.On("Set", mock.Anything, &testOrders[2]).
					Return(nil).
					Once()
			},
			workers: 1,
		},
		{
			name: "empty orders list",
			setupMocks: func(repo *MockOrderProvider, _ *MockOrderSetter) {
				repo.On("GetAll", mock.Anything).
					Return([]domain.Order{}, nil).
					Once()
			},
			workers: 2,
		},
		{
			name: "repository error",
			setupMocks: func(repo *MockOrderProvider, _ *MockOrderSetter) {
				repo.On("GetAll", mock.Anything).
					Return(nil, errors.New("database error")).
					Once()
			},
			workers:       2,
			expectedError: "failed to get all orders: database error",
		},
		{
			name: "cache set error",
			setupMocks: func(repo *MockOrderProvider, cache *MockOrderSetter) {
				repo.On("GetAll", mock.Anything).
					Return(testOrders, nil).
					Once()
				cache.On("Set", mock.Anything, &testOrders[0]).
					Return(nil).
					Once()
				cache.On("Set", mock.Anything, &testOrders[1]).
					Return(errors.New("cache error")).
					Once()
				cache.On("Set", mock.Anything, &testOrders[2]).
					Return(nil).
					Maybe()
			},
			workers:       2,
			expectedError: "failed to restore orders: cache error",
		},
		{
			name: "context cancellation",
			setupMocks: func(repo *MockOrderProvider, cache *MockOrderSetter) {
				repo.On("GetAll", mock.Anything).
					Return(testOrders, nil).
					Once()
				cache.On("Set", mock.Anything, mock.Anything).
					Return(context.Canceled).
					Maybe()
			},
			workers:       2,
			expectedError: "failed to restore orders: context canceled",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockOrderProvider(t)
			mockCache := NewMockOrderSetter(t)
			tt.setupMocks(mockRepo, mockCache)

			ctx := context.Background()
			err := Restore(ctx, mockRepo, mockCache, tt.workers)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestRestore_Concurrency(t *testing.T) {
	t.Parallel()

	testOrders := make([]domain.Order, 10)
	for i := range testOrders {
		testOrders[i] = domain.Order{OrderUID: string(rune('a' + i))}
	}

	mockRepo := NewMockOrderProvider(t)
	mockCache := NewMockOrderSetter(t)

	mockRepo.On("GetAll", mock.Anything).
		Return(testOrders, nil).
		Once()

	for i := range testOrders {
		mockCache.On("Set", mock.Anything, &testOrders[i]).
			Return(nil).
			Once()
	}

	ctx := context.Background()
	err := Restore(ctx, mockRepo, mockCache, 5)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestRestore_ZeroWorkers(t *testing.T) {
	t.Parallel()

	mockRepo := NewMockOrderProvider(t)
	mockCache := NewMockOrderSetter(t)

	testOrders := []domain.Order{{OrderUID: "test"}}
	mockRepo.On("GetAll", mock.Anything).
		Return(testOrders, nil).
		Once()
	mockCache.On("Set", mock.Anything, &testOrders[0]).
		Return(nil).
		Once()

	ctx := context.Background()
	err := Restore(ctx, mockRepo, mockCache, 0)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestRestore_ContextPropagation(t *testing.T) {
	t.Parallel()

	mockRepo := NewMockOrderProvider(t)
	mockCache := NewMockOrderSetter(t)

	testOrders := []domain.Order{{OrderUID: "test"}}
	mockRepo.On("GetAll", mock.Anything).
		Return(testOrders, nil).
		Once()
	mockCache.On("Set", mock.Anything, &testOrders[0]).
		Return(nil).
		Once()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Restore(ctx, mockRepo, mockCache, 1)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestRestore_NilOrders(t *testing.T) {
	t.Parallel()

	mockRepo := NewMockOrderProvider(t)
	mockCache := NewMockOrderSetter(t)

	var emptyOrders []domain.Order
	mockRepo.On("GetAll", mock.Anything).
		Return(emptyOrders, nil).
		Once()

	ctx := context.Background()
	err := Restore(ctx, mockRepo, mockCache, 2)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
