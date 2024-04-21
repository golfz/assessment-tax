package tax

import (
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
			name:      "rate 0%: income=100,000 deduction.personal=0; expect tax=0",
			info:      TaxInformation{TotalIncome: 100_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 0.0},
		},
		{
			name:      "rate 0%: income=150,000 deduction.personal=0; expect tax=0",
			info:      TaxInformation{TotalIncome: 150_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 0.0},
		},
		{
			name:      "rate 10%: income=150,001 deduction.personal=0; expect tax=35,000",
			info:      TaxInformation{TotalIncome: 150_001.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 0.1},
		},
		{
			name:      "rate 10%: income=500,000 deduction.personal=60,000; expect tax=29,000 (EXP01)",
			info:      TaxInformation{TotalIncome: 500_000.0},
			deduction: Deduction{Personal: 60_000.0},
			want:      TaxResult{Tax: 29_000.0},
		},
		{
			name:      "rate 10%: income=500,000 deduction.personal=0; expect tax=35,000",
			info:      TaxInformation{TotalIncome: 500_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 35_000.0},
		},
		{
			name:      "rate 10% + 15%: income=500,001 deduction.personal=0; expect tax=35,000.15",
			info:      TaxInformation{TotalIncome: 500_001.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 35_000.15},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := CalculateTax(tc.info, tc.deduction)

			// Assert
			if err != nil {
				t.Errorf("CalculateTax() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("CalculateTax() = %v, want %v", got, tc.want)
			}
		})
	}
}
