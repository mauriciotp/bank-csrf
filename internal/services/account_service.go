package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mauriciotp/bank-csrf/internal/store/pgstore"
	"github.com/mauriciotp/bank-csrf/internal/usecase/account"
)

var (
	ErrDuplicatedEmail   = errors.New("email already in use")
	ErrInvalidAccountID  = errors.New("invalid account id")
	ErrAccountNotFound   = errors.New("account not found")
	ErrSameAccount       = errors.New("cannot use same account to transfer")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type TransferRow struct {
	SenderID         uuid.UUID      `json:"sender_id"`
	RecipientID      uuid.UUID      `json:"recipient"`
	SenderBalance    pgtype.Numeric `json:"sender_balance"`
	RecipientBalance pgtype.Numeric `json:"recipient_balance"`
}

type AccountsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewAccountsService(pool *pgxpool.Pool) *AccountsService {
	return &AccountsService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (s *AccountsService) CreateAccount(ctx context.Context, req account.CreateAccountReq) (pgstore.Account, error) {
	income, err := toNumeric(req.Income)
	if err != nil {
		return pgstore.Account{}, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return pgstore.Account{}, err
	}

	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	params := pgstore.CreateAccountParams{
		MobilePhone: pgtype.Text{
			String: req.MobilePhone,
			Valid:  req.MobilePhone != "",
		},
	}

	switch req.Type {
	case account.TypeNaturalPerson:
		person, err := qtx.CreateNaturalPerson(ctx, pgstore.CreateNaturalPersonParams{
			FullName:      req.Name,
			Age:           req.Age,
			Email:         req.Email,
			MonthlyIncome: income,
		})
		if err != nil {
			return pgstore.Account{}, mapUniqueViolation(err)
		}
		params.NaturalPersonID = pgtype.UUID{
			Bytes: person.ID,
			Valid: true,
		}
	case account.TypeLegalEntity:
		person, err := qtx.CreateLegalEntity(ctx, pgstore.CreateLegalEntityParams{
			TradeName:      req.Name,
			Age:            req.Age,
			CorporateEmail: req.Email,
			Revenue:        income,
		})
		if err != nil {
			return pgstore.Account{}, mapUniqueViolation(err)
		}
		params.LegalEntityID = pgtype.UUID{
			Bytes: person.ID,
			Valid: true,
		}
	}

	category, err := qtx.GetCategoryByName(ctx, req.Category())
	if err != nil {
		return pgstore.Account{}, err
	}
	params.CategoryID = category.ID

	acc, err := qtx.CreateAccount(ctx, params)
	if err != nil {
		return pgstore.Account{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return pgstore.Account{}, err
	}

	return acc, nil
}

func (s *AccountsService) GetAccountBalance(ctx context.Context, id string) (pgstore.GetAccountBalanceRow, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return pgstore.GetAccountBalanceRow{}, ErrInvalidAccountID
	}

	balanceRow, err := s.queries.GetAccountBalance(ctx, parsedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.GetAccountBalanceRow{}, ErrAccountNotFound
		}
		return pgstore.GetAccountBalanceRow{}, err
	}

	return balanceRow, nil
}

func (s *AccountsService) Deposit(ctx context.Context, id string, req account.DepositReq) (pgstore.DepositRow, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return pgstore.DepositRow{}, ErrInvalidAccountID
	}

	amount, err := toNumeric(req.Amount)
	if err != nil {
		return pgstore.DepositRow{}, err
	}

	params := pgstore.DepositParams{
		ID:     parsedId,
		Amount: amount,
	}

	depositRow, err := s.queries.Deposit(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.DepositRow{}, ErrAccountNotFound
		}
		return pgstore.DepositRow{}, err
	}

	return depositRow, nil
}

func (s *AccountsService) Withdraw(ctx context.Context, id string, req account.WithdrawReq) (pgstore.WithdrawRow, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return pgstore.WithdrawRow{}, ErrInvalidAccountID
	}

	amount, err := toNumeric(req.Amount)
	if err != nil {
		return pgstore.WithdrawRow{}, err
	}

	params := pgstore.WithdrawParams{
		ID:     parsedId,
		Amount: amount,
	}

	withdrawRow, err := s.queries.Withdraw(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.WithdrawRow{}, ErrAccountNotFound
		}
		return pgstore.WithdrawRow{}, err
	}

	return withdrawRow, nil
}

func (s *AccountsService) Transfer(ctx context.Context, req account.TransferReq) (TransferRow, error) {
	parsedSenderId, err := uuid.Parse(req.SenderID)
	if err != nil {
		return TransferRow{}, ErrInvalidAccountID
	}

	parsedRecipientId, err := uuid.Parse(req.RecipientID)
	if err != nil {
		return TransferRow{}, ErrInvalidAccountID
	}

	if parsedSenderId == parsedRecipientId {
		return TransferRow{}, ErrSameAccount
	}

	amount, err := toNumeric(req.Amount)
	if err != nil {
		return TransferRow{}, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return TransferRow{}, err
	}

	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	accountsIds := []uuid.UUID{parsedSenderId, parsedRecipientId}
	ids, err := qtx.LockAccounts(ctx, accountsIds)
	if err != nil {
		return TransferRow{}, err
	}
	if len(ids) != 2 {
		return TransferRow{}, ErrAccountNotFound
	}

	withdrawRow, err := qtx.Withdraw(ctx, pgstore.WithdrawParams{
		Amount: amount,
		ID:     parsedSenderId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TransferRow{}, ErrInsufficientFunds
		}
		return TransferRow{}, err
	}

	depositRow, err := qtx.Deposit(ctx, pgstore.DepositParams{
		Amount: amount,
		ID:     parsedRecipientId,
	})
	if err != nil {
		return TransferRow{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return TransferRow{}, err
	}

	return TransferRow{
		SenderID:         parsedSenderId,
		RecipientID:      parsedRecipientId,
		SenderBalance:    withdrawRow.Balance,
		RecipientBalance: depositRow.Balance,
	}, nil
}

func (s *AccountsService) CloseAccount(ctx context.Context, id string) (pgstore.CloseAccountRow, error) {
	accountId, err := uuid.Parse(id)
	if err != nil {
		return pgstore.CloseAccountRow{}, ErrInvalidAccountID
	}

	closeAccountRow, err := s.queries.CloseAccount(ctx, accountId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.CloseAccountRow{}, ErrAccountNotFound
		}
		return pgstore.CloseAccountRow{}, err
	}

	return closeAccountRow, nil
}

func mapUniqueViolation(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrDuplicatedEmail
	}
	return err
}

func toNumeric(v float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(strconv.FormatFloat(v, 'f', 2, 64))
	return n, err
}
