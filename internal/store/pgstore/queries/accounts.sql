-- name: CreateNaturalPerson :one
INSERT INTO natural_persons (
  full_name,
  age,
  email,
  monthly_income
) VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: CreateLegalEntity :one
INSERT INTO legal_entities (
  trade_name,
  age,
  corporate_email,
  revenue
) VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;
