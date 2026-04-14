package postgres

import (
	"context"

	"charts-user-service/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	query := `
	INSERT INTO users (id, name, email, password_hash, about)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING created_at
	`

	err := r.db.QueryRow(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.About,
	).Scan(&user.CreatedAt)

	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {

	query := `
		SELECT id, name, email, password_hash, about, created_at
		FROM users
		WHERE email = $1
	`

	u := &user.User{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.About,
		&u.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, name, email, password_hash, about, created_at
		FROM users
		WHERE id = $1
	`

	user := &user.User{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.About,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *user.User) error {
	query := `
		UPDATE users
		SET name = $2,
		    email = $3,
		    password_hash = $4,
		    about = $5
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.About,
	)

	return err
}
