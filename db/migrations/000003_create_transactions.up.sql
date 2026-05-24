CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    from_wallet_id UUID REFERENCES wallets (id),
    to_wallet_id UUID REFERENCES wallets (id),
    amount NUMERIC(19, 4) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

CREATE INDEX idx_transactions_from_wallet ON transactions (from_wallet_id);
CREATE INDEX idx_transactions_to_wallet ON transactions (to_wallet_id);
