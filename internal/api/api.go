package api

import "github.com/go-chi/chi/v5"

type API struct {
	Router *chi.Mux
}
