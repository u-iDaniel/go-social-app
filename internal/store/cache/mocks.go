package cache

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/u-iDaniel/go-social-app/internal/store"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	// Mock implementation
	args := m.Called(userID)
	return nil, args.Error(1)
}

func (m *MockUserStore) Set(ctx context.Context, user *store.User) error {
	// Mock implementation
	args := m.Called(user)
	return args.Error(0)
}
