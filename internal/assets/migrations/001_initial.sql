-- +migrate Up
CREATE TABLE usdcTransactions (
    id SERIAL PRIMARY KEY,
    tx_hash VARCHAR(66) NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    value NUMERIC(38, 18) NOT NULL,
    timestamp BIGINT NOT NULL,
    block_number BIGINT NOT NULL
);

-- +migrate Down
DROP TABLE usdcTransactions;
