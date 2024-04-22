//go:build unit

package tax

import (
	"reflect"
	"testing"
)

func TestValidateDeduction(t *testing.T) {
	defaultDeduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}
	// Arrange
	testcases := []struct {
		name      string
		deduction Deduction
		want      error
	}{
		// Default deduction
		{
			name: "default deduction; expect no error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			want: nil,
		},
		// personal deduction
		{
			name: "personal deduction < min; expect error",
			deduction: Deduction{
				Personal: ConstraintMinPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			want: ErrInvalidDeduction,
		},
		{
			name: "personal deduction = min; expect no error",
			deduction: Deduction{
				Personal: ConstraintMinPersonalDeduction + 0.1,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			want: nil,
		},
		{
			name: "personal deduction = max; expect no error",
			deduction: Deduction{
				Personal: ConstraintMaxPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			want: nil,
		},
		{
			name: "personal deduction > max; expect error",
			deduction: Deduction{
				Personal: ConstraintMaxPersonalDeduction + 0.1,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			want: ErrInvalidDeduction,
		},
		// KReceipt deduction
		{
			name: "KReceipt deduction < min; expect error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMinKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
			want: ErrInvalidDeduction,
		},
		{
			name: "KReceipt deduction = min; expect no error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMinKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
			want: nil,
		},
		{
			name: "KReceipt deduction = max; expect no error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMaxKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
			want: nil,
		},
		{
			name: "KReceipt deduction > max; expect error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMaxKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
			want: ErrInvalidDeduction,
		},
		// Donation deduction
		{
			name: "Donation deduction > max; expect error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: ConstraintMaxDonationDeduction + 0.1,
			},
			want: ErrInvalidDeduction,
		},
		{
			name: "Donation deduction = max; expect no error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: ConstraintMaxDonationDeduction,
			},
			want: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := validateDeduction(tc.deduction)

			// Assert
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("validateDeduction() = %v, want %v", got, tc.want)
			}
		})
	}
}

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
			name:      "rate 15%: income=500,001 deduction.personal=0; expect tax=35,000.15",
			info:      TaxInformation{TotalIncome: 500_001.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 35_000.15},
		},
		{
			name:      "rate 15%: income=750,000 deduction.personal=0; expect tax=72,500",
			info:      TaxInformation{TotalIncome: 750_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 72_500.0},
		},
		{
			name:      "rate 15%: income=1,000,000 deduction.personal=0; expect tax=110,000",
			info:      TaxInformation{TotalIncome: 1_000_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 110_000.0},
		},
		{
			name:      "rate 20%: income=1,000,001 deduction.personal=0; expect tax=110,000.20",
			info:      TaxInformation{TotalIncome: 1_000_001.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 110_000.20},
		},
		{
			name:      "rate 20%: income=1,500,000 deduction.personal=0; expect tax=210,000",
			info:      TaxInformation{TotalIncome: 1_500_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 210_000.0},
		},
		{
			name:      "rate 20%: income=2,000,000 deduction.personal=0; expect tax=310,000",
			info:      TaxInformation{TotalIncome: 2_000_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 310_000.0},
		},
		{
			name:      "rate 35%: income=2,000,001 deduction.personal=0; expect tax=310,000.35",
			info:      TaxInformation{TotalIncome: 2_000_001.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 310_000.35},
		},
		{
			name:      "rate 35%: income=10,000,000 deduction.personal=0; expect tax=3,110,000",
			info:      TaxInformation{TotalIncome: 10_000_000.0},
			deduction: Deduction{Personal: 0.0},
			want:      TaxResult{Tax: 3_110_000.0},
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
