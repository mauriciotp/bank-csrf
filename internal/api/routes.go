package api

import "github.com/go-chi/chi/v5"

func (a *API) BindRoutes() {
	a.Router.Route("/api", func(r chi.Router) {
	})
}
