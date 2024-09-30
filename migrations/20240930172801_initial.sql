-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.contracts (
    id INTEGER PRIMARY KEY NOT NULL,
    is_active BOOLEAN NOT NULL,
    agent_token VARCHAR(254),
    patient_name VARCHAR(254),
    patient_email VARCHAR(254),
    locale VARCHAR(5) NULL,
    libre_client INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.contracts;
-- +goose StatementEnd
