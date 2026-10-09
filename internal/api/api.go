package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/mauriciotp/bank-csrf/internal/services"
)

type API struct {
	Router          *chi.Mux
	AccountsService *services.AccountsService
}
