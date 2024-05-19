-- migrate:up
CREATE TABLE loans (
    LIKE template_table INCLUDING ALL,
    borrower_id INT NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    interest_rate NUMERIC(5, 2) NOT NULL,
    duration_weeks INT NOT NULL,
    outstanding NUMERIC(15, 2) NOT NULL,
    delinquent_weeks INT NOT NULL,
    CONSTRAINT fk_borrower
        FOREIGN KEY(borrower_id) 
        REFERENCES borrowers(id)
);

-- migrate:down
DROP TABLE loans;