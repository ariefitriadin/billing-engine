package usecase

import (
	"billing-engine/internal/domain"
	"context"
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the LoanRepository
type MockLoanRepository struct {
	mock.Mock
}

// CreateBillingSchedule implements domain.LoanRepository.
func (m *MockLoanRepository) CreateBillingSchedule(ctx context.Context, schedule *domain.BillingSchedule) error {
	panic("unimplemented")
}

// CreateLoan implements domain.LoanRepository.
func (m *MockLoanRepository) CreateLoan(ctx context.Context, borrowerID uint, loan *domain.Loan) (uint, error) {
	panic("unimplemented")
}

// GetBillingSchedule implements domain.LoanRepository.
func (m *MockLoanRepository) GetBillingSchedule(ctx context.Context, loanId uint) (*domain.BillingSchedule, error) {
	panic("unimplemented")
}

// GetLoansWithBorrower implements domain.LoanRepository.
func (m *MockLoanRepository) GetLoansWithBorrower(ctx context.Context, limit uint, offset uint) ([]domain.LoanWithBorrower, error) {
	panic("unimplemented")
}

// IsDelinquent implements domain.LoanRepository.
func (m *MockLoanRepository) IsDelinquent(ctx context.Context, loanID uint) (*domain.CheckDelinquentAmount, error) {
	panic("unimplemented")
}

// UpdateBillingSchedule implements domain.LoanRepository.
func (m *MockLoanRepository) UpdateBillingSchedule(ctx context.Context, schedule *domain.BillingSchedule) error {
	panic("unimplemented")
}

// UpdateLoan implements domain.LoanRepository.
func (m *MockLoanRepository) UpdateLoan(ctx context.Context, loan *domain.Loan, schedule *domain.BillingSchedule) error {
	panic("unimplemented")
}

// UpdateRepaymentSchedule implements domain.LoanRepository.
func (m *MockLoanRepository) UpdateRepaymentSchedule(ctx context.Context, loan *domain.Loan) error {
	panic("unimplemented")
}

func (m *MockLoanRepository) GetLoanByID(ctx context.Context, loanID uint) (*domain.Loan, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(*domain.Loan), args.Error(1)
}

func TestGetOutstanding(t *testing.T) {
	mockRepo := new(MockLoanRepository)
	loanUsecase := NewLoanUsecase(mockRepo)
	ctx := context.Background()
	loanID := uint(1)

	// Setting up the expected values
	expectedLoan := &domain.Loan{
		Outstanding: pgtype.Numeric{Int: big.NewInt(500), Exp: 0},
	}
	// Assuming the value 500.00 needs to be set
	// expectedLoan.Outstanding.Set(500.00)

	mockRepo.On("GetLoanByID", ctx, loanID).Return(expectedLoan, nil)

	// Running the test
	result, err := loanUsecase.GetOutstanding(ctx, loanID)

	// Asserting the results
	assert.NoError(t, err)
	assert.Equal(t, 500.00, result)
	mockRepo.AssertExpectations(t)
}
