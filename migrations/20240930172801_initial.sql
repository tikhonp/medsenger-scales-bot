-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.contracts (
    id INTEGER PRIMARY KEY NOT NULL,
    is_active BOOLEAN NOT NULL,
    agent_token VARCHAR(254),
    patient_name VARCHAR(254),
    patient_email VARCHAR(254),
    locale VARCHAR(5) NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.contracts;
-- +goose StatementEnd
