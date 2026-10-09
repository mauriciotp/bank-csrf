package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mauriciotp/bank-csrf/internal/jsonutils"
	"github.com/mauriciotp/bank-csrf/internal/services"
	"github.com/mauriciotp/bank-csrf/internal/usecase/account"
)

func (a *API) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeJSON[account.CreateAccountReq](w, r)
	if err != nil {
		if problems != nil {
			_ = jsonutils.EncodeJSON(w, http.StatusUnprocessableEntity, problems)
			return
		}
		_ = jsonutils.EncodeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	acc, err := a.AccountsService.CreateAccount(r.Context(), data)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmail) {
			_ = jsonutils.EncodeError(w, http.StatusConflict, err.Error())
			return
		}
		slog.Error("failed to create account", "error", err)
		_ = jsonutils.EncodeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	_ = jsonutils.EncodeJSON(w, http.StatusCreated, acc)
}

func (a *API) handleGetAccountBalance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	balanceRow, err := a.AccountsService.GetAccountBalance(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAccountID):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrAccountNotFound):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		default:
			_ = jsonutils.EncodeError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
	}

	_ = jsonutils.EncodeJSON(w, http.StatusOK, balanceRow)
}

func (a *API) handleDeposit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, problems, err := jsonutils.DecodeJSON[account.DepositReq](w, r)
	if err != nil {
		if len(problems) > 0 {
			_ = jsonutils.EncodeJSON(w, http.StatusUnprocessableEntity, problems)
			return
		}
		_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
		return
	}

	depositRow, err := a.AccountsService.Deposit(r.Context(), id, data)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAccountID):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrAccountNotFound):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		default:
			_ = jsonutils.EncodeError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
	}

	_ = jsonutils.EncodeJSON(w, http.StatusOK, depositRow)
}

func (a *API) handleWithdraw(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, problems, err := jsonutils.DecodeJSON[account.WithdrawReq](w, r)
	if err != nil {
		if len(problems) > 0 {
			_ = jsonutils.EncodeJSON(w, http.StatusUnprocessableEntity, problems)
			return
		}
		_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
		return
	}

	withdrawRow, err := a.AccountsService.Withdraw(r.Context(), id, data)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAccountID):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrAccountNotFound):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		default:
			_ = jsonutils.EncodeError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
	}

	_ = jsonutils.EncodeJSON(w, http.StatusOK, withdrawRow)
}

func (a *API) handleTransfer(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeJSON[account.TransferReq](w, r)
	if err != nil {
		if len(problems) > 0 {
			_ = jsonutils.EncodeJSON(w, http.StatusUnprocessableEntity, problems)
			return
		}
		_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
		return
	}

	transferRow, err := a.AccountsService.Transfer(r.Context(), data)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAccountID):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrAccountNotFound):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrSameAccount):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrInsufficientFunds):
			_ = jsonutils.EncodeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		default:
			_ = jsonutils.EncodeError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
	}

	_ = jsonutils.EncodeJSON(w, http.StatusOK, transferRow)
}

func (a *API) handleCloseAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	closeAccountRow, err := a.AccountsService.CloseAccount(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAccountID):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrAccountNotFound):
			_ = jsonutils.EncodeError(w, http.StatusBadRequest, err.Error())
			return
		default:
			_ = jsonutils.EncodeError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
	}

	_ = jsonutils.EncodeJSON(w, http.StatusOK, closeAccountRow)
}
