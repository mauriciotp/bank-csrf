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
) RETURNING *;
