-- Reverses 000004_payments_tokenize.up.sql. Restores the legacy schema but
-- cannot recover dropped data: card_no was never preserved (PCI-driven
-- destruction in the up migration), and amount is reconstructed from
-- amount_cents only when that column still exists.

ALTER TABLE payments
    ADD COLUMN IF NOT EXISTS card_no BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS amount  REAL   NOT NULL DEFAULT 0;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM   information_schema.columns
        WHERE  table_name  = 'payments'
        AND    column_name = 'amount_cents'
    ) THEN
        EXECUTE 'UPDATE payments
                 SET    amount = (amount_cents::numeric / 100)::real
                 WHERE  amount_cents IS NOT NULL';
    END IF;
END $$;

ALTER TABLE payments DROP COLUMN IF EXISTS card_token;
ALTER TABLE payments DROP COLUMN IF EXISTS amount_cents;
ALTER TABLE payments DROP COLUMN IF EXISTS currency_code;

ALTER TABLE payments ALTER COLUMN card_no DROP DEFAULT;
ALTER TABLE payments ALTER COLUMN amount  DROP DEFAULT;
