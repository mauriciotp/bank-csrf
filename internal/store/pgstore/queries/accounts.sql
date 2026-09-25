-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE category_name = $1;

-- name: CreateAccount :one
INSERT INTO accounts (
  legal_entity_id,
  natural_person_id,
  mobile_phone,
  category_id
) VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING *;

-- name: GetAccountBalance :one
SELECT
  id,
  balance
FROM accounts
WHERE id = $1 AND closed_at IS NULL;

-- name: Deposit :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id) AND closed_at IS NULL
RETURNING id, balance;

-- name: Withdraw :one
UPDATE accounts
SET balance = balance - sqlc.arg(amount)
WHERE
  id = sqlc.arg(id)
  AND closed_at IS NULL
  AND balance >= sqlc.arg(amount)
RETURNING id, balance;

-- name: LockAccounts :many
SELECT id
FROM accounts
WHERE id = ANY(sqlc.arg(ids)::uuid[])
ORDER BY id
FOR UPDATE;

-- name: CloseAccount :one
UPDATE accounts
SET closed_at = NOW()
WHERE id = $1 AND closed_at IS NULL AND balance = 0
RETURNING id, closed_at;
