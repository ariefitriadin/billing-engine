package domain

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Loan struct {
	ID                uint
	Amount            pgtype.Numeric
	InterestRate      pgtype.Numeric
	DurationWeeks     int
	Outstanding       pgtype.Numeric
	DelinquentWeeks   int
	InstallmentAmount pgtype.Numeric
}

type LoanRepository interface {
	GetLoanByID(ctx context.Context, loanID uint) (*Loan, error)
	UpdateLoan(ctx context.Context, loan *Loan, schedule *BillingSchedule) error
	GetLoansWithBorrower(ctx context.Context, limit, offset uint) ([]LoanWithBorrower, error)
	CreateLoan(ctx context.Context, borrowerID uint, loan *Loan) (uint, error)
	CreateBillingSchedule(ctx context.Context, schedule *BillingSchedule) error
	UpdateBillingSchedule(ctx context.Context, schedule *BillingSchedule) error
	GetBillingSchedule(ctx context.Context, loanId uint) (*BillingSchedule, error)
	IsDelinquent(ctx context.Context, loanID uint) (*CheckDelinquentAmount, error)
	UpdateRepaymentSchedule(ctx context.Context, loan *Loan) error
}

type LoanWithBorrower struct {
	LoanID        uint
	BorrowerID    uint
	BorrowerName  string
	Amount        pgtype.Numeric
	InterestRate  pgtype.Numeric
	DurationWeeks int
	Outstanding   pgtype.Numeric
}

type BillingSchedule struct {
	ID      uint
	LoanID  uint
	Week    uint
	Amount  pgtype.Numeric
	DueDate pgtype.Date
	Paid    pgtype.Bool
}

type CheckDelinquentAmount struct {
	LoanID       uint
	TotalWeek    int
	Amount       int64
	IsDelinquent bool
}
