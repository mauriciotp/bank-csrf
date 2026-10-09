package account

import (
	"context"

	"github.com/mauriciotp/bank-csrf/internal/validator"
)

type TransferReq struct {
	SenderID    string  `json:"sender_id"`
	RecipientID string  `json:"recipient_id"`
	Amount      float64 `json:"amount"`
}

func (req TransferReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Min(req.Amount, 1), "amount", "cannot be negative")
	eval.CheckField(validator.IsValidUUID(req.SenderID), "sender_id", "invalid uuid")
	eval.CheckField(validator.IsValidUUID(req.RecipientID), "recipient_id", "invalid uuid")

	return eval
}
