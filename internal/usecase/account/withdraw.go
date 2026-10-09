package account

import (
	"context"

	"github.com/mauriciotp/bank-csrf/internal/validator"
)

type WithdrawReq struct {
	Amount float64 `json:"amount"`
}

func (req WithdrawReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Min(req.Amount, 1), "amount", "cannot be negative")

	return eval
}
