package util

import (
	"fmt"
	"testing"
)

func TestRound(t *testing.T) {
	tests := []float64{
		12.7845,
		-24.741245,
		7841541154.451278,
		-784545134.45754,
	}
	for _, flt := range tests {
		fmt.Printf("%f \n", Round(flt, 2))
	}
}

func TestAbs(t *testing.T) {
	test_case := []interface{}{
		int(-2),
		int8(-1),
		float32(-0.618),
	}
	for _, item := range test_case {
		fmt.Printf("%v\n", Abs(item))
	}
}
