CREATE TABLE goods (
    id         BIGSERIAL PRIMARY KEY,
    order_id   BIGINT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);