-- +migrate Up
CREATE TABLE usdcTransactions (
    id SERIAL PRIMARY KEY,
    txHash VARCHAR(66) NOT NULL,
    fromAddress VARCHAR(42) NOT NULL,
    toAddress VARCHAR(42) NOT NULL,
    value NUMERIC(38, 18) NOT NULL,
    timestamp BIGINT NOT NULL,

    UNIQUE (txHash, fromAddress, toAddress, value, timestamp)
);
-- +migrate Down
DROP TABLE usdcTransactions;
