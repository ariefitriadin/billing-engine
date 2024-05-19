package utils

import (
	"fmt"
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func BigFloatToNumeric(f *big.Float) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	if f == nil {
		return numeric, fmt.Errorf("input big.Float is nil")
	}
	// Convert big.Float to string to use the Set method of pgtype.Numeric
	str := f.Text('f', -1)
	err := numeric.Scan(str)
	if err != nil {
		return numeric, err
	}
	return numeric, nil
}

func NumericToBigFloat(n pgtype.Numeric) (*big.Float, error) {

	intnumber := n.Int.Int64()
	expnumber := n.Exp

	return big.NewFloat(float64(intnumber) * math.Pow(10, float64(expnumber))), nil
}

func NumericToFloat64(n pgtype.Numeric) (float64, error) {
	intnumber := n.Int.Int64()
	expnumber := n.Exp

	return float64(intnumber) * math.Pow(10, float64(expnumber)), nil
}

func Float64ToNumeric(f float64) (pgtype.Numeric, error) {
	bigFloat := big.NewFloat(f)
	return BigFloatToNumeric(bigFloat)
}

// CeilBigFloat returns the smallest integer value greater than or equal to x.
func CeilBigFloat(x *big.Float) *big.Float {
	if x == nil {
		return nil
	}

	// Create a new big.Int to hold the integer part
	intPart := new(big.Int)
	// Create a new big.Float to hold the remainder
	fracPart := new(big.Float)

	// Get the integer part and the fractional part
	x.Int(intPart)
	fracPart.Sub(x, new(big.Float).SetInt(intPart))

	// If there's a fractional part and the value is positive, add 1 to the integer part
	if fracPart.Cmp(big.NewFloat(0)) > 0 {
		intPart.Add(intPart, big.NewInt(1))
	}

	// Set the result to the integer part
	result := new(big.Float).SetInt(intPart)

	return result
}
