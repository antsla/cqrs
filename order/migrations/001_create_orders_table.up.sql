CREATE TABLE orders (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);