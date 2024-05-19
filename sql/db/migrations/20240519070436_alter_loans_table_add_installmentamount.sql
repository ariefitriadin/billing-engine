-- migrate:up
ALTER TABLE loans
ADD COLUMN installment_amount NUMERIC(15, 2);

-- migrate:down
ALTER TABLE loans
DROP COLUMN installment_amount;
