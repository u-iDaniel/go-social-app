package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type RolesStore struct {
	db *sql.DB
}

func (s *RolesStore) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `
        SELECT id, name, description, level
        FROM roles
        WHERE name = $1;
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := &Role{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		name,
	).Scan(&role.ID, &role.Name, &role.Description, &role.Level)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return role, nil
}
