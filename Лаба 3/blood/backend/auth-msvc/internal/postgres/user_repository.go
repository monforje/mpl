package postgres

import (
	"auth/internal/core"
	"auth/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	pg *Postgres
}

func NewUserRepository(pg *Postgres) core.UserRepository {
	return &UserRepository{pg: pg}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, name, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pg.Pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("не удалось создать пользователя: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, name, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	user := &model.User{}
	err := r.pg.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось получить пользователя: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, email, name, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	user := &model.User{}
	err := r.pg.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось получить пользователя: %w", err)
	}
	return user, nil
}
