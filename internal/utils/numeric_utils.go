/*
Package utils provides utility functions for converting between big.Float and pgtype.Numeric types, and also includes functions to convert pgtype.Numeric to float64 and vice versa. Additionally, it contains a function to compute the ceiling of a big.Float.

Example Usage:
bigFloat := big.NewFloat(123.456)
numeric, err := BigFloatToNumeric(bigFloat)

newBigFloat, err := NumericToBigFloat(numeric)

floatVal, err := NumericToFloat64(numeric)

numericFromFloat, err := Float64ToNumeric(123.456)

ceilValue := CeilBigFloat(big.NewFloat(123.456))

Code Analysis:
Inputs:
- big.Float: A floating-point number with arbitrary precision.
- pgtype.Numeric: A PostgreSQL numeric type used to handle arbitrary precision numbers.
- float64: A floating-point number with double precision.

Flow:
1. BigFloatToNumeric converts a big.Float to pgtype.Numeric by first converting the big.Float to a string and then scanning it into a pgtype.Numeric.
2. NumericToBigFloat and NumericToFloat64 convert a pgtype.Numeric to big.Float and float64 respectively, using the integer and exponent parts of pgtype.Numeric.
3. Float64ToNumeric converts a float64 to pgtype.Numeric by first converting the float64 to big.Float and then using BigFloatToNumeric.
4. CeilBigFloat calculates the ceiling of a big.Float by separating it into integer and fractional parts, and adjusting the integer part if there is a non-zero fractional component.

Outputs:
- pgtype.Numeric: The converted numeric type from big.Float or float64.
- big.Float: The converted or adjusted big floating-point number.
- float64: The converted double precision floating-point number.
- error: Potential errors that can occur during conversions or calculations.
*/
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

	intPart, fracPart := new(big.Int), new(big.Float)
	x.Int(intPart)
	fracPart.Sub(x, new(big.Float).SetInt(intPart))

	if fracPart.Cmp(big.NewFloat(0)) > 0 {
		intPart.Add(intPart, big.NewInt(1))
	}

	return new(big.Float).SetInt(intPart)
}
