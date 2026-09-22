CREATE TABLE natural_persons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  full_name VARCHAR(255) NOT NULL,
  age INTEGER NOT NULL,
  monthly_income DECIMAL(15, 2) NOT NULL DEFAULT 0,
  email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE legal_entities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  trade_name VARCHAR(255) NOT NULL,
  age INTEGER NOT NULL,
  revenue DECIMAL(15, 2) NOT NULL DEFAULT 0,
  corporate_email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  category_name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  natural_person_id UUID REFERENCES natural_persons (id),
  legal_entity_id UUID REFERENCES legal_entities (id),
  category_id INTEGER NOT NULL REFERENCES categories (id),
  mobile_phone VARCHAR(20),
  balance DECIMAL(15, 2) NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  closed_at TIMESTAMPTZ,
  CONSTRAINT chk_unique_owner CHECK (
    (natural_person_id IS NOT NULL AND legal_entity_id IS NULL)
    OR (natural_person_id IS NULL AND legal_entity_id IS NOT NULL)
  )
);

---- create above / drop below ----
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS legal_entities;
DROP TABLE IF EXISTS natural_persons;
