-- name: CreateAccount :one
INSERT INTO accounts (
  legal_entity_id,
  natural_person_id,
  balance,
  mobile_phone,
  category_id
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
) RETURNING *;

-- name: GetBalanceByAccount :one
SELECT balance FROM accounts
WHERE id = $1;

-- name: UpdateBalance :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING balance;

-- name: CloseAccount :exec
DELETE FROM accounts
WHERE id = $1;
