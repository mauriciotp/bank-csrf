package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mauriciotp/bank-csrf/internal/store/pgstore"
)

type API struct {
	Router  *chi.Mux
	Pool    *pgxpool.Pool
	Queries *pgstore.Queries
}
