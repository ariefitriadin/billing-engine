package usecase

import (
	"billing-engine/internal/domain"
	"billing-engine/internal/utils"
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

type LoanUsecase interface {
	GetOutstanding(ctx context.Context, loanID uint) (float64, error)
	IsDelinquent(ctx context.Context, loanID uint) (*domain.CheckDelinquentAmount, error)
	MakePayment(ctx context.Context, loanID uint, amount float64) error
	GetLoansWithBorrower(ctx context.Context, limit, offset uint) ([]domain.LoanWithBorrower, error)
	CreateLoan(ctx context.Context, borrowerID uint, amount float64, interestRate, durationWeeks int) (uint, error)
}

type loanUsecase struct {
	loanRepo domain.LoanRepository
}

func NewLoanUsecase(lr domain.LoanRepository) LoanUsecase {
	return &loanUsecase{loanRepo: lr}
}

func (lu *loanUsecase) GetOutstanding(ctx context.Context, loanID uint) (float64, error) {
	loan, err := lu.loanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		return 0, err
	}
	outstandingFloat, err := utils.NumericToFloat64(loan.Outstanding)
	if err != nil {
		return 0, err
	}
	return outstandingFloat, nil
}

func (lu *loanUsecase) IsDelinquent(ctx context.Context, loanID uint) (*domain.CheckDelinquentAmount, error) {
	check, err := lu.loanRepo.IsDelinquent(ctx, loanID)
	if err != nil {
		return nil, err
	}
	check.IsDelinquent = check.TotalWeek >= 2
	return check, nil
}

func (lu *loanUsecase) MakePayment(ctx context.Context, loanID uint, amount float64) error {

	loan, err := lu.loanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		return fmt.Errorf("failed to get loan by id %d: %w", loanID, err)
	}
	// Convert pgtype.Numeric to big.Float for comparison
	amountFloat := new(big.Float).SetFloat64(amount)

	outstandingFloat, err := utils.NumericToBigFloat(loan.Outstanding)
	if err != nil {
		return errors.New("NumericToBigFloat outstanding amount")
	}

	if outstandingFloat.Cmp(big.NewFloat(0)) == 0 {
		return errors.New("loan is already full paid")
	}

	installmentAmountFloat, err := utils.NumericToBigFloat(loan.InstallmentAmount)
	if err != nil {
		return errors.New("NumericToBigFloat for installment amount")
	}

	checkDelinquentAmount, err := lu.loanRepo.IsDelinquent(ctx, loanID)
	if err != nil {
		return errors.New("failed to check delinquent amount")
	}

	if checkDelinquentAmount.IsDelinquent && amount != float64(checkDelinquentAmount.Amount) {
		return fmt.Errorf("please repay the arrears amount first, for: %.2f", float64(checkDelinquentAmount.Amount))
	}

	if amountFloat.Cmp(installmentAmountFloat) != 0 && !checkDelinquentAmount.IsDelinquent {
		return fmt.Errorf("payment amount not equal to installment amount: %.2f", installmentAmountFloat)
	}

	outstandingBalance := new(big.Float).Sub(outstandingFloat, amountFloat)
	loan.Outstanding, err = utils.BigFloatToNumeric(outstandingBalance)
	if err != nil {
		return errors.New("BigFloatToNumeric for outstanding balance")
	}

	if checkDelinquentAmount.IsDelinquent {
		return lu.loanRepo.UpdateRepaymentSchedule(ctx, loan)
	}

	nearestBillingSchedule, err := lu.loanRepo.GetBillingSchedule(ctx, loanID)
	if err != nil {
		return errors.New("failed to get billing schedule")
	}
	nearestBillingSchedule.Paid = pgtype.Bool{Bool: true, Valid: true}

	return lu.loanRepo.UpdateLoan(ctx, loan, nearestBillingSchedule)
}

func (lu *loanUsecase) GetLoansWithBorrower(ctx context.Context, limit, offset uint) ([]domain.LoanWithBorrower, error) {
	return lu.loanRepo.GetLoansWithBorrower(ctx, limit, offset)
}

func (lu *loanUsecase) CreateLoan(ctx context.Context, borrowerID uint, amount float64, interestRate, durationWeeks int) (uint, error) {

	rate := float64(interestRate) / 100
	interestAmount := amount * rate
	outstanding := new(big.Float).SetFloat64(amount + interestAmount)

	outstandingNumeric, err := utils.BigFloatToNumeric(outstanding)
	if err != nil {
		return 0, fmt.Errorf("failed to convert outstanding amount: %w", err)
	}

	amountNumeric, err := utils.BigFloatToNumeric(new(big.Float).SetFloat64(amount))
	if err != nil {
		return 0, fmt.Errorf("failed to convert amount: %w", err)
	}

	interestRateNumeric, err := utils.BigFloatToNumeric(new(big.Float).SetFloat64(rate))
	if err != nil {
		return 0, fmt.Errorf("failed to convert interest rate: %w", err)
	}

	weeklyAmount := utils.CeilBigFloat(new(big.Float).Quo(outstanding, big.NewFloat(float64(durationWeeks))))
	weeklyAmountNumeric, err := utils.BigFloatToNumeric(weeklyAmount)
	if err != nil {
		return 0, fmt.Errorf("failed to convert weekly amount: %w", err)
	}

	loan := &domain.Loan{
		Amount:            amountNumeric,
		InterestRate:      interestRateNumeric,
		DurationWeeks:     durationWeeks,
		Outstanding:       outstandingNumeric,
		InstallmentAmount: weeklyAmountNumeric,
	}

	loanID, err := lu.loanRepo.CreateLoan(ctx, borrowerID, loan)
	if err != nil {
		return 0, err
	}

	return loanID, nil
}

func (lu *loanUsecase) UpdateBillingSchedule(ctx context.Context, schedule *domain.BillingSchedule) error {
	return lu.loanRepo.UpdateBillingSchedule(ctx, schedule)
}
