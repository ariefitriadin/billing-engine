-- migrate:up
CREATE TABLE billing_schedule (
    id SERIAL PRIMARY KEY,
    loan_id INT NOT NULL,
    week INT NOT NULL,
    amount NUMERIC NOT NULL,
    due_date DATE NOT NULL,
    FOREIGN KEY (loan_id) REFERENCES loans(id)
);

-- migrate:down
DROP TABLE billing_schedule;
