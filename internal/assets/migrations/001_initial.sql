-- +migrate Up
CREATE TABLE usdcTransactions (
    txHash VARCHAR(66) PRIMARY KEY,
    fromAddress VARCHAR(42) NOT NULL,
    toAddress VARCHAR(42) NOT NULL,
    value NUMERIC(25, 18) NOT NULL,
    timestamp BIGINT NOT NULL
);

-- +migrate Down
DROP TABLE usdcTransactions;
