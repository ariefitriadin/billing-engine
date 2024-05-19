package http

import (
	"billing-engine/internal/usecase"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type LoanHandler struct {
	lu usecase.LoanUsecase
}

func NewLoanHandler(e *echo.Echo, lu usecase.LoanUsecase) {
	handler := &LoanHandler{lu: lu}
	e.GET("/loans/:id/outstanding", handler.GetOutstanding)
	e.GET("/loans/:id/delinquent", handler.IsDelinquent)
	e.POST("/loans/:id/payment", handler.MakePayment)
	e.GET("/loans", handler.GetLoansWithBorrower)
	e.POST("/loans", handler.CreateLoan)
}

// @Summary Get outstanding amount
// @Description Get the current outstanding amount for a loan
// @ID get-outstanding
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]float64
// @Router /loans/{id}/outstanding [get]
func (lh *LoanHandler) GetOutstanding(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid loan ID"})
	}
	outstanding, err := lh.lu.GetOutstanding(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]float64{"outstanding": outstanding})
}

// @Summary Check if loan is delinquent
// @Description Check if the loan is delinquent
// @ID is-delinquent
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]bool
// @Router /loans/{id}/delinquent [get]
func (lh *LoanHandler) IsDelinquent(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid loan ID"})
	}
	delinquent, err := lh.lu.IsDelinquent(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, delinquent)
}

// @Summary Make a payment
// @Description Make a payment on the loan
// @ID make-payment
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param amount body float64 true "Payment Amount"
// @Success 200 {object} map[string]string
// @Router /loans/{id}/payment [post]
func (lh *LoanHandler) MakePayment(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid loan ID"})
	}
	var request struct {
		Amount float64 `json:"amount"`
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := lh.lu.MakePayment(ctx, uint(id), request.Amount); err != nil {
		log.Printf("Error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "payment successful"})
}

// @Summary Get loans with borrower information
// @Description Get a list of loans with borrower information
// @ID get-loans-with-borrower
// @Produce json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Success 200 {array} domain.LoanWithBorrower
// @Router /loans [get]
func (lh *LoanHandler) GetLoansWithBorrower(c echo.Context) error {
	ctx := c.Request().Context()
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid limit"})
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid offset"})
	}

	loans, err := lh.lu.GetLoansWithBorrower(ctx, uint(limit), uint(offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, loans)
}

// @Summary Create a new loan
// @Description Create a new loan and generate a billing schedule
// @ID create-loan
// @Accept json
// @Produce json
// @Param borrower_id body int true "Borrower ID"
// @Param amount body float64 true "Loan Amount"
// @Param interest_rate body float64 true "Interest Rate"
// @Param duration_weeks body int true "Duration in Weeks"
// @Success 200 {object} map[string]uint
// @Router /loans [post]
func (lh *LoanHandler) CreateLoan(c echo.Context) error {
	var request struct {
		BorrowerID    uint    `json:"borrower_id"`
		Amount        float64 `json:"amount"`
		InterestRate  int     `json:"interest_rate"`
		DurationWeeks int     `json:"duration_weeks"`
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()
	loanID, err := lh.lu.CreateLoan(ctx, request.BorrowerID, request.Amount, request.InterestRate, request.DurationWeeks)
	if err != nil {
		if err == context.DeadlineExceeded {
			return c.JSON(http.StatusRequestTimeout, map[string]string{"error": "request timed out"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]uint{"loan_id": loanID})
}
