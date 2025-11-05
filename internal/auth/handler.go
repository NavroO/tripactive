package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/logout", h.Logout)
	return r
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
}
