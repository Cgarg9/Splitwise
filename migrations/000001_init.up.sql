-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Enums
DO $$ BEGIN
  CREATE TYPE connection_status AS ENUM ('Pending', 'Accepted', 'Blocked');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
  CREATE TYPE group_role AS ENUM ('Member', 'Admin');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Users
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE,
  email VARCHAR(100) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  phone_number VARCHAR(15) UNIQUE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Connections
CREATE TABLE IF NOT EXISTS connections (
    first_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    second_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    status connection_status NOT NULL DEFAULT 'Pending',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (first_user_id, second_user_id),
    CHECK (first_user_id <> second_user_id),
    UNIQUE (second_user_id, first_user_id)
);

-- Groups 
CREATE TABLE IF NOT EXISTS groups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(100) NOT NULL,
  description TEXT,
  default_currency CHAR(10) NOT NULL DEFAULT 'INR',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  created_by UUID REFERENCES users(id),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Group Members
CREATE TABLE IF NOT EXISTS group_members (
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role group_role NOT NULL DEFAULT 'Member',
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (group_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members(user_id);

-- Expenses
CREATE TABLE IF NOT EXISTS expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE DEFAULT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    description TEXT,
    currency CHAR(10) NOT NULL DEFAULT 'INR',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    expense_date TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) NOT NULL,
    expense_type VARCHAR(50) NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS expenses_group_date_idx
  ON expenses(group_id, expense_date DESC);

-- Expense Split 
CREATE TABLE IF NOT EXISTS expense_splits (
    expense_id UUID REFERENCES expenses(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    paid_amount Numeric(10, 2) NOT NULL DEFAULT 0 CHECK (paid_amount >= 0),
    owed_amount DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (owed_amount >= 0),
    PRIMARY KEY (expense_id, user_id)
);