package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUsersStore{},
	}
}

type MockUsersStore struct {
	// users []User // In-memory slice to store users (unused for now)
}

func (m *MockUsersStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	// Mock implementation
	return nil
}

func (m *MockUsersStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	// Mock implementation
	return nil
}

func (m *MockUsersStore) GetByID(ctx context.Context, id int64) (*User, error) {
	// Mock implementation
	return &User{ID: id}, nil
}

func (m *MockUsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	// Mock implementation
	return &User{}, nil
}

func (m *MockUsersStore) Activate(ctx context.Context, token string) error {
	// Mock implementation
	return nil
}

func (m *MockUsersStore) Delete(ctx context.Context, id int64) error {
	// Mock implementation
	return nil
}

func (m *MockUsersStore) Follow(ctx context.Context, followerID, userID int64) error {
	// Mock implementation
	return nil
}

func (m *MockUsersStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	// Mock implementation
	return nil
}
