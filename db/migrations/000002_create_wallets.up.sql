CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL DEFAULT 'My Wallet',
    balance NUMERIC(19, 4) NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

CREATE OR REPLACE FUNCTION check_wallet_limit()
RETURNS TRIGGER AS $$
BEGIN
    IF (SELECT COUNT(*) FROM wallets WHERE user_id = NEW.user_id) >= 3 THEN
        RAISE EXCEPTION 'Maximum wallet limit reached';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_wallet_limit
    BEFORE INSERT ON wallets
    FOR EACH ROW EXECUTE FUNCTION check_wallet_limit();

CREATE INDEX idx_wallets_user_id ON wallets (user_id);
