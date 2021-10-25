CREATE TABLE orders (
    id         BIGINT PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);