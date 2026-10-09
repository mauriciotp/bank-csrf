package account

import (
	"context"

	"github.com/mauriciotp/bank-csrf/internal/validator"
)

const (
	TypeNaturalPerson = "natural_person"
	TypeLegalEntity   = "legal_entity"
)

type CreateAccountReq struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Age         int32   `json:"age"`
	Email       string  `json:"email"`
	Income      float64 `json:"income"`
	MobilePhone string  `json:"mobile_phone"`
}

func (req CreateAccountReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(
		validator.PermittedValue(req.Type, TypeNaturalPerson, TypeLegalEntity),
		"type", "must be 'natural_person' or 'legal_entity'",
	)

	eval.CheckField(validator.NotBlank(req.Name), "name", "name is required")
	eval.CheckField(validator.MaxChars(req.Name, 255), "name", "must have at most 255 characters")

	if req.Type == TypeNaturalPerson {
		eval.CheckField(validator.Min(req.Age, 18), "age", "must be at least 18")
	} else {
		eval.CheckField(validator.Min(req.Age, 0), "age", "cannot be negative")
	}

	eval.CheckField(validator.IsEmail(req.Email), "email", "must be a valid email")
	eval.CheckField(validator.MaxChars(req.Email, 255), "email", "must have at most 255 characters")

	eval.CheckField(validator.Min(req.Income, 0), "income", "cannot be negative")

	if req.MobilePhone != "" {
		eval.CheckField(validator.Matches(req.MobilePhone, validator.PhoneRX), "mobile_phone", "must have 10 to 15 digits")
	}

	return eval
}

func (req CreateAccountReq) Category() string {
	if req.Type == TypeNaturalPerson {
		switch {
		case req.Income < 5_000:
			return "standard"
		case req.Income < 20_000:
			return "premium"
		default:
			return "private"
		}
	}

	switch {
	case req.Income < 360_000:
		return "micro_business"
	case req.Income < 4_800_000:
		return "small_business"
	default:
		return "corporate"
	}
}
