package api

import "github.com/go-chi/chi/v5"

func (a *API) BindRoutes() {
	a.Router.Route("/accounts", func(r chi.Router) {
		r.Post("/", a.handleCreateAccount)
		r.Get("/{id}/balance", a.handleGetAccountBalance)
		r.Post("/{id}/deposit", a.handleDeposit)
		r.Post("/{id}/withdraw", a.handleWithdraw)
		r.Post("/transfer", a.handleTransfer)
		r.Delete("/{id}", a.handleCloseAccount)
	})
}
