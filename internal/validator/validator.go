package validator

import "context"

type Validator interface {
	Valid(ctx context.Context) Evaluator
}

type Evaluator map[string]string

func (e *Evaluator) AddFieldError(key, message string) {
	if *e == nil {
		*e = make(map[string]string)
	}

	if _, exists := (*e)[key]; !exists {
		(*e)[key] = message
	}
}

func (e *Evaluator) CheckField(valid bool, key, message string) {
	if !valid {
		e.AddFieldError(key, message)
	}
}
