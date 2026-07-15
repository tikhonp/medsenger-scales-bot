-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.contracts
ADD COLUMN patient_sex VARCHAR(255) NULL,
ADD COLUMN patient_age INTEGER NULL,
ADD COLUMN patient_height FLOAT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.contracts
DROP COLUMN patient_sex,
DROP COLUMN patient_age,
DROP COLUMN patient_height;
-- +goose StatementEnd
