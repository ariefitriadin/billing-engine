-- migrate:up
ALTER TABLE billing_schedule
ADD COLUMN paid BOOLEAN DEFAULT FALSE;

-- migrate:down
ALTER TABLE billing_schedule
DROP COLUMN paid;