package auth

import "context"

type Service interface {
	Register(ctx context.Context) ([]Ad, error)
	Login(ctx context.Context, body CreateAdRequest) (Ad, error)
	Refresh(ctx context.Context, body CreateAdRequest) (Ad, error)
	Logout(ctx context.Context, body CreateAdRequest) (Ad, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
