-- Tokenize payments: replace raw PAN/REAL amount with an opaque card_token,
-- integer cents, and an ISO-4217 currency code. PCI-DSS compliance requires
-- that full card numbers not be stored at rest here.
--
-- Order matters: add new columns first, backfill from the legacy columns
-- inside the same transaction the migrate runner gives us, then drop the old
-- columns. This preserves the `amount` data on existing rows. The historical
-- `card_no` is intentionally not preserved -- by design we do not have a way
-- to retroactively tokenize it, and keeping it would defeat the purpose of
-- this migration.

ALTER TABLE payments
    ADD COLUMN IF NOT EXISTS card_token    TEXT   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS amount_cents  BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS currency_code TEXT   NOT NULL DEFAULT 'USD';

-- Backfill amount_cents from the legacy REAL `amount` column, casting through
-- numeric to avoid binary-float rounding drift. Wrapped in a DO block so the
-- migration is re-runnable: if `amount` was already dropped on a prior run,
-- this no-ops instead of erroring.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM   information_schema.columns
        WHERE  table_name  = 'payments'
        AND    column_name = 'amount'
    ) THEN
        EXECUTE 'UPDATE payments
                 SET    amount_cents = ROUND(amount::numeric * 100)::bigint
                 WHERE  amount IS NOT NULL';
    END IF;
END $$;

ALTER TABLE payments DROP COLUMN IF EXISTS card_no;
ALTER TABLE payments DROP COLUMN IF EXISTS amount;

-- Drop the temporary defaults so future inserts must specify these fields.
ALTER TABLE payments ALTER COLUMN card_token   DROP DEFAULT;
ALTER TABLE payments ALTER COLUMN amount_cents DROP DEFAULT;
