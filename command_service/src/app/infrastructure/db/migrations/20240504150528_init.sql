-- +goose Up
-- +goose StatementBegin
CREATE TABLE command
(
    id SERIAL PRIMARY KEY ,
    description VARCHAR(255),
    user_id BIGINT NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TYPE STATUS as ENUM ('processing', 'success', 'failed');
CREATE TABLE command_outbox
(
    event_uuid UUID PRIMARY KEY,
    command_id INTEGER references command(id),
    status STATUS
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE command_outbox;
DROP TYPE STATUS;
DROP TABLE command;
-- +goose StatementEnd
