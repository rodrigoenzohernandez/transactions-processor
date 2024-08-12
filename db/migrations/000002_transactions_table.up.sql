CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS "dev";

CREATE TABLE IF NOT EXISTS "dev".transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    date VARCHAR(255) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL
);
