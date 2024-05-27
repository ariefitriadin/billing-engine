package utils

import (
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

// BigFloatToNumeric converts a non-nil big.Float to a pgtype.Numeric successfully
func TestBigFloatToNumericSuccess(t *testing.T) {
	bigFloat := big.NewFloat(123.456)
	expectedNumeric := pgtype.Numeric{Int: big.NewInt(123456), Exp: -3}
	numeric, err := BigFloatToNumeric(bigFloat)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if numeric.Int.Cmp(expectedNumeric.Int) != 0 || numeric.Exp != expectedNumeric.Exp {
		t.Errorf("Expected %v, got %v", expectedNumeric, numeric)
	}
}

// BigFloatToNumeric returns an error if the input big.Float is nil
func TestBigFloatToNumericNilInput(t *testing.T) {
	_, err := BigFloatToNumeric(nil)
	if err == nil {
		t.Error("Expected an error for nil big.Float input, but got none")
	}
}
