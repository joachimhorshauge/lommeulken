-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
                         "id" integer PRIMARY KEY,
                         "email" varchar UNIQUE NOT NULL,
                         "username" varchar UNIQUE NOT NULL,
                         "password_hash" varchar,
                         "salt" varchar,
                         "verification_hash" varchar,
                         "role" varchar,
                         "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
                         "modified_at" timestamp DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "users";
-- +goose StatementEnd