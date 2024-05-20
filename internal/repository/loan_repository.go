package repository

import (
	"billing-engine/internal/domain"
	"billing-engine/sql/billingengine"
	"context"
	"fmt"
	"time"

	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type loanRepository struct {
	queries *billingengine.Queries
	db      *pgxpool.Pool
}

type noRowsError struct {
	sqlNoRow string
	noRows   string
}

func (e noRowsError) ErrSqlNoRow() string {
	e.sqlNoRow = "sql: no rows in result set"
	return e.sqlNoRow
}

func (e noRowsError) ErrNoRow() string {
	e.noRows = "no rows in result set"
	return e.noRows
}

func NewLoanRepository(db *pgxpool.Pool) domain.LoanRepository {
	return &loanRepository{queries: billingengine.New(db), db: db}
}

func (r *loanRepository) IsDelinquent(ctx context.Context, loanID uint) (*domain.CheckDelinquentAmount, error) {
	check, err := r.queries.CheckDelinquentAmount(ctx, int32(loanID))
	ernorow := noRowsError{}
	if err != nil {
		if err.Error() == ernorow.ErrSqlNoRow() || err.Error() == ernorow.ErrNoRow() {
			return &domain.CheckDelinquentAmount{
				LoanID:       uint(loanID),
				IsDelinquent: false,
				TotalWeek:    0,
				Amount:       0,
			}, nil
		} else {
			return nil, fmt.Errorf("loan repo: failed to check delinquent amount: %w", err)
		}
	}
	return &domain.CheckDelinquentAmount{
		LoanID:       uint(check.LoanID),
		TotalWeek:    int(check.TotalWeek),
		Amount:       check.Sum,
		IsDelinquent: check.TotalWeek >= 2,
	}, nil
}

func (r *loanRepository) GetLoanByID(ctx context.Context, loanID uint) (*domain.Loan, error) {
	loan, err := r.queries.GetLoanByID(ctx, int32(loanID))
	if err != nil {
		log.Printf("failed to get loan by id: %v", err)
		return nil, err
	}
	return &domain.Loan{
		ID:                uint(loan.ID),
		Amount:            loan.Amount,
		InterestRate:      loan.InterestRate,
		DurationWeeks:     int(loan.DurationWeeks),
		Outstanding:       loan.Outstanding,
		DelinquentWeeks:   int(loan.DelinquentWeeks),
		InstallmentAmount: loan.InstallmentAmount,
	}, nil
}

func (r *loanRepository) UpdateLoan(ctx context.Context, loan *domain.Loan, schedule *domain.BillingSchedule) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	params := billingengine.UpdateLoanParams{
		Amount:          loan.Amount,
		InterestRate:    loan.InterestRate,
		DurationWeeks:   int32(loan.DurationWeeks),
		Outstanding:     loan.Outstanding,
		DelinquentWeeks: int32(loan.DelinquentWeeks),
		ID:              int32(loan.ID),
	}

	err = r.queries.WithTx(tx).UpdateLoan(ctx, params)
	if err != nil {
		log.Printf("failed to update loan: %v", err)
		return fmt.Errorf("failed to update loan: %w", err)
	}

	err = r.queries.WithTx(tx).UpdateBillingSchedule(ctx, billingengine.UpdateBillingScheduleParams{
		Paid:   schedule.Paid,
		LoanID: int32(schedule.LoanID),
		Week:   int32(schedule.Week),
	})
	if err != nil {
		log.Printf("failed to update billing schedule: %v", err)
		return fmt.Errorf("failed to update billing schedule: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit UpdateLoan transaction: %w", err)
	}
	return nil
}

func (r *loanRepository) GetLoansWithBorrower(ctx context.Context, limit, offset uint) ([]domain.LoanWithBorrower, error) {
	params := billingengine.GetLoansWithBorrowerParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	loans, err := r.queries.GetLoansWithBorrower(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get loans with borrower: %w", err)
	}

	var result []domain.LoanWithBorrower
	for _, loan := range loans {
		result = append(result, domain.LoanWithBorrower{
			LoanID:        uint(loan.LoanID),
			BorrowerID:    uint(loan.BorrowerID),
			BorrowerName:  loan.BorrowerName,
			Amount:        loan.Amount,
			InterestRate:  loan.InterestRate,
			DurationWeeks: int(loan.DurationWeeks),
			Outstanding:   loan.Outstanding,
		})
	}
	return result, nil
}

func (r *loanRepository) CreateLoan(ctx context.Context, borrowerID uint, loan *domain.Loan) (uint, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	loanID, err := r.queries.WithTx(tx).CreateLoan(ctx, billingengine.CreateLoanParams{
		BorrowerID:        int32(borrowerID),
		Amount:            loan.Amount,
		InterestRate:      loan.InterestRate,
		DurationWeeks:     int32(loan.DurationWeeks),
		Outstanding:       loan.Outstanding,
		DelinquentWeeks:   int32(0),
		InstallmentAmount: loan.InstallmentAmount,
	})
	if err != nil {
		log.Printf("failed to create loan: %v", err)
		return 0, fmt.Errorf("failed to create loan: %w", err)
	}
	var billingSchedules billingengine.CreateBillingSchedulesParams
	for week := 1; week <= loan.DurationWeeks; week++ {
		dueDate := time.Now().AddDate(0, 0, 7*week)
		billingSchedules.Column1 = append(billingSchedules.Column1, int32(loanID))
		billingSchedules.Column2 = append(billingSchedules.Column2, int32(week))
		billingSchedules.Column3 = append(billingSchedules.Column3, loan.InstallmentAmount)
		billingSchedules.Column4 = append(billingSchedules.Column4, pgtype.Date{Time: dueDate, Valid: true})
		billingSchedules.Column5 = append(billingSchedules.Column5, false)
	}

	err = r.queries.WithTx(tx).CreateBillingSchedules(ctx, billingSchedules)
	if err != nil {
		log.Printf("failed to create billing schedules: %v", err)
		return 0, fmt.Errorf("failed to create billing schedules: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to commit CreateLoan transaction: %w", err)
	}

	return uint(loanID), nil
}

func (r *loanRepository) CreateBillingSchedule(ctx context.Context, schedule *domain.BillingSchedule) error {

	err := r.queries.CreateBillingSchedule(ctx, billingengine.CreateBillingScheduleParams{
		LoanID:  int32(schedule.LoanID),
		Week:    int32(schedule.Week),
		Amount:  schedule.Amount,
		DueDate: schedule.DueDate,
		Paid:    schedule.Paid,
	})
	if err != nil {
		log.Printf("failed to create billing schedule: %v", err)
		return fmt.Errorf("failed to create billing schedule: %w", err)
	}
	return nil
}

func (r *loanRepository) GetBillingSchedule(ctx context.Context, loanID uint) (*domain.BillingSchedule, error) {

	billSchedule, err := r.queries.GetBillingSchedule(ctx, int32(loanID))

	if err != nil {
		return nil, fmt.Errorf("failed to get billing schedule: %w", err)
	}
	return &domain.BillingSchedule{
		ID:      uint(billSchedule.ID),
		LoanID:  uint(billSchedule.LoanID),
		Week:    uint(billSchedule.Week),
		Amount:  billSchedule.Amount,
		DueDate: billSchedule.DueDate,
		Paid:    billSchedule.Paid,
	}, nil
}

func (r *loanRepository) UpdateBillingSchedule(ctx context.Context, schedule *domain.BillingSchedule) error {
	err := r.queries.UpdateBillingSchedule(ctx, billingengine.UpdateBillingScheduleParams{
		Paid:   schedule.Paid,
		LoanID: int32(schedule.LoanID),
		Week:   int32(schedule.Week),
	})
	if err != nil {
		log.Printf("failed to update billing schedule: %v", err)
		return fmt.Errorf("failed to update billing schedule: %w", err)
	}
	return nil
}

func (r *loanRepository) UpdateRepaymentSchedule(ctx context.Context, loan *domain.Loan) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin UpdateRepaymentSchedule transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	params := billingengine.UpdateLoanParams{
		Amount:        loan.Amount,
		InterestRate:  loan.InterestRate,
		DurationWeeks: int32(loan.DurationWeeks),
		Outstanding:   loan.Outstanding,
		ID:            int32(loan.ID),
	}

	err = r.queries.WithTx(tx).UpdateLoan(ctx, params)
	if err != nil {
		log.Printf("failed to update loan: %v", err)
		return fmt.Errorf("failed to update loan: %w", err)
	}

	err = r.queries.WithTx(tx).UpdateRepaymentSchedule(ctx, int32(loan.ID))
	if err != nil {
		return fmt.Errorf("failed to update repayment schedule: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit UpdateRepaymentSchedule transaction: %w", err)
	}

	return nil
}
