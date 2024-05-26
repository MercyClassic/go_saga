-- +goose Up
-- +goose StatementBegin
CREATE TABLE account
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(255),
    username VARCHAR(255) UNIQUE NOT NULL,
    balance  NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0)
);
CREATE TABLE account_balance_log
(
    event_id   UUID PRIMARY KEY,
    amount     NUMERIC(10, 2) NOT NULL,
    user_id    INTEGER        NOT NULL REFERENCES account (id),
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE account_balance_log;
DROP TABLE account;
-- +goose StatementEnd
