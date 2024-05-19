-- name: GetLoanByID :one
SELECT id, amount, interest_rate, duration_weeks, outstanding, delinquent_weeks, installment_amount
FROM loans
WHERE id = $1;

-- name: UpdateLoan :exec
UPDATE loans
SET amount = $1, interest_rate = $2, duration_weeks = $3, outstanding = $4, delinquent_weeks = $5
WHERE id = $6;

-- name: GetLoansByBorrowerID :many
SELECT id, amount, interest_rate, duration_weeks, outstanding, delinquent_weeks
FROM loans
WHERE borrower_id = $1;

-- name: GetLoansWithBorrower :many
SELECT 
    loans.id AS loan_id, 
    borrowers.id AS borrower_id, 
    borrowers.name AS borrower_name, 
    loans.amount, 
    loans.interest_rate, 
    loans.duration_weeks, 
    loans.outstanding, 
    loans.delinquent_weeks,
    loans.installment_amount
FROM loans
JOIN borrowers ON loans.borrower_id = borrowers.id
LIMIT $1 OFFSET $2;


-- name: CreateLoan :one
INSERT INTO loans (borrower_id, amount, interest_rate, duration_weeks, outstanding, delinquent_weeks, installment_amount)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: CreateBillingSchedule :exec
INSERT INTO billing_schedule (loan_id, week, amount, due_date, paid)
VALUES ($1, $2, $3, $4, $5);

-- name: CreateBillingSchedules :exec
INSERT INTO billing_schedule (loan_id, week, amount, due_date, paid)
VALUES (
    unnest($1::int[]),
    unnest($2::int[]),
    unnest($3::numeric[]),
    unnest($4::date[]),
    unnest($5::boolean[])
);


-- name: UpdateBillingSchedule :exec
UPDATE billing_schedule
SET paid = $1
WHERE loan_id = $2 AND week = $3;

-- name: GetBillingSchedule :one
SELECT id, loan_id, week, amount, due_date, paid
FROM billing_schedule
WHERE loan_id = $1 AND paid = false ORDER BY week LIMIT 1;

-- name: CheckDelinquentAmount :one
SELECT loan_id, count(1) as total_week, sum(amount)
FROM billing_schedule
WHERE loan_id = $1 AND paid = false AND due_date <  now()
GROUP BY loan_id;

-- name: UpdateRepaymentSchedule :exec
UPDATE billing_schedule
SET paid = true
WHERE loan_id = $1 AND due_date < now();
