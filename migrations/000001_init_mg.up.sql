BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(50) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL,
    password VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS user_detail (
    id VARCHAR(50) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    firstname VARCHAR(100),
    lastname VARCHAR(100),
    gender VARCHAR(6),
    city VARCHAR(100),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS user_preference (
    id VARCHAR(50) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    preferred_gender VARCHAR(6),
    preferred_city VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS payment (
    id VARCHAR(50) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    invoice_number VARCHAR(100) NOT NULL,
    bank_name VARCHAR(50),
    payment_status VARCHAR(20),
    amount_due INTEGER NOT NULL,
    va_number INTEGER NOT NULL,
    token VARCHAR(200) NOT NULL,
    request_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    payment_date TIMESTAMPTZ NOT NULL,
    deadline TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS subscription (
    id VARCHAR(50) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    payment_id VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL
);

CREATE TABLE IF NOT EXISTS swipe (
    id VARCHAR(50) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    swipe_user_id VARCHAR(50) NOT NULL,
    is_like BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL
);

COMMIT;