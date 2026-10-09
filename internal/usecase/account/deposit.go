package account

import (
	"context"

	"github.com/mauriciotp/bank-csrf/internal/validator"
)

type DepositReq struct {
	Amount float64 `json:"amount"`
}

func (req DepositReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Min(req.Amount, 0), "amount", "cannot be negative")

	return eval
}
