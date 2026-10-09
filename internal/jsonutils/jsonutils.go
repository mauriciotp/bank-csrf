package jsonutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mauriciotp/bank-csrf/internal/validator"
)

func EncodeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

func DecodeJSON[T validator.Validator](w http.ResponseWriter, r *http.Request) (T, map[string]string, error) {
	var data T

	r.Body = http.MaxBytesReader(w, r.Body, 1<<10)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&data); err != nil {
		switch {
		case errors.Is(err, io.EOF):
			return data, nil, errors.New("request body cannot be empty")
		case errors.Is(err, io.ErrUnexpectedEOF):
			return data, nil, errors.New("request body contains malformed JSON")
		default:
			return data, nil, err
		}
	}

	if problems := data.Valid(r.Context()); len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T: %d problems", data, len(problems))
	}

	return data, nil, nil
}

func EncodeError(w http.ResponseWriter, status int, msg string) error {
	return EncodeJSON(w, status, map[string]string{"error": msg})
}
