package tax

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCalculateTax_onlyTotalIncome(t *testing.T) {
	// Arrange
	testcases := []struct {
		name      string
		info      TaxInformation
		deduction Deduction
		want      TaxResult
	}{
		{
			name:      "tax level 1: income=100,000 deduction.personal=0; expect tax=0",
			info:      TaxInformation{TotalIncome: 100_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 0.0},
		},
		{
			name:      "tax level 2: income=500,000 deduction.personal=60,000; expect tax=29,000",
			info:      TaxInformation{TotalIncome: 500_000.0},
			deduction: Deduction{Personal: 60_000.0},
			want:      TaxResult{Tax: 29_000.0},
		},
		{
			name:      "tax level 2: income=500,000 deduction.personal=0; expect tax=35,000",
			info:      TaxInformation{TotalIncome: 500_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 35_000.0},
		},
	}

	for _, tc := range testcases {
		// Act
		got, err := CalculateTax(tc.info, tc.deduction)

		// Assert
		fmt.Printf("got: %#v\n", got)
		if err != nil {
			t.Errorf("CalculateTax() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("CalculateTax() = %v, want %v", got, tc.want)
		}
	}
}
