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
) RETURNING *;
