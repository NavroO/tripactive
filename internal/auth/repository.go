package auth

import (
	"context"
	"database/sql"
)

type Repository interface {
	Register(ctx context.Context) ([]Ad, error)
	Login(ctx context.Context, body CreateAdRequest) (Ad, error)
	Refresh(ctx context.Context, body CreateAdRequest) (Ad, error)
	Logout(ctx context.Context, body CreateAdRequest) (Ad, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
