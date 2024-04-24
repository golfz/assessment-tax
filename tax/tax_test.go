//go:build unit

package tax

import (
	"github.com/stretchr/testify/assert"
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
		wantError error
	}{
		// Default deduction
		{
			name: "default deduction; expect no error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantError: nil,
		},
		// personal deduction
		{
			name: "personal deduction < min; expect error",
			deduction: Deduction{
				Personal: ConstraintMinPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantError: ErrInvalidDeduction,
		},
		{
			name: "personal deduction = min; expect no error",
			deduction: Deduction{
				Personal: ConstraintMinPersonalDeduction + 0.1,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantError: nil,
		},
		{
			name: "personal deduction = max; expect no error",
			deduction: Deduction{
				Personal: ConstraintMaxPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantError: nil,
		},
		{
			name: "personal deduction > max; expect error",
			deduction: Deduction{
				Personal: ConstraintMaxPersonalDeduction + 0.1,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantError: ErrInvalidDeduction,
		},
		// KReceipt deduction
		{
			name: "KReceipt deduction < min; expect error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMinKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
			wantError: ErrInvalidDeduction,
		},
		{
			name: "KReceipt deduction = min; expect no error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMinKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
			wantError: nil,
		},
		{
			name: "KReceipt deduction = max; expect no error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMaxKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
			wantError: nil,
		},
		{
			name: "KReceipt deduction > max; expect error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMaxKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
			wantError: ErrInvalidDeduction,
		},
		// Donation deduction
		{
			name: "Donation deduction > max; expect error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: ConstraintMaxDonationDeduction + 0.1,
			},
			wantError: ErrInvalidDeduction,
		},
		{
			name: "Donation deduction = max; expect no error",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: ConstraintMaxDonationDeduction,
			},
			wantError: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotError := validateDeduction(tc.deduction)

			// Assert
			assert.Equal(t, tc.wantError, gotError)
		})
	}
}

func TestCalculateTax_Success_InputOnlyTotalIncome(t *testing.T) {
	// Arrange
	defaultDeduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}
	testcases := []struct {
		name      string
		info      TaxInformation
		deduction Deduction
		want      TaxResult
	}{
		{
			name:      "rate 0%: income=100,000 deduction.personal=0; expect tax=0",
			info:      TaxInformation{TotalIncome: 100_000.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 0.0},
		},
		{
			name:      "rate 0%: income=150,000 deduction.personal=0; expect tax=0",
			info:      TaxInformation{TotalIncome: 150_000.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 0.0},
		},
		{
			name:      "rate 10%: income=150,001 deduction.personal=0; expect tax=35,000",
			info:      TaxInformation{TotalIncome: 150_001.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 0.1},
		},
		{
			name:      "rate 10%: income=500,000 deduction.personal=60,000; expect tax=29,000 (EXP01)",
			info:      TaxInformation{TotalIncome: 500_000.0},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 29_000.0},
		},
		{
			name:      "rate 10%: income=500,000 deduction.personal=0; expect tax=35,000",
			info:      TaxInformation{TotalIncome: 500_000.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 35_000.0},
		},
		{
			name:      "rate 15%: income=500,001 deduction.personal=0; expect tax=35,000.15",
			info:      TaxInformation{TotalIncome: 500_001.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 35_000.15},
		},
		{
			name:      "rate 15%: income=750,000 deduction.personal=0; expect tax=72,500",
			info:      TaxInformation{TotalIncome: 750_000.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 72_500.0},
		},
		{
			name:      "rate 15%: income=1,000,000 deduction.personal=0; expect tax=110,000",
			info:      TaxInformation{TotalIncome: 1_000_000.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 110_000.0},
		},
		{
			name:      "rate 20%: income=1,000,001 deduction.personal=0; expect tax=110,000.20",
			info:      TaxInformation{TotalIncome: 1_000_001.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 110_000.20},
		},
		{
			name:      "rate 20%: income=1,500,000 deduction.personal=0; expect tax=210,000",
			info:      TaxInformation{TotalIncome: 1_500_000.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 210_000.0},
		},
		{
			name:      "rate 20%: income=2,000,000 deduction.personal=0; expect tax=310,000",
			info:      TaxInformation{TotalIncome: 2_000_000.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 310_000.0},
		},
		{
			name:      "rate 35%: income=2,000,001 deduction.personal=0; expect tax=310,000.35",
			info:      TaxInformation{TotalIncome: 2_000_001.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 310_000.35},
		},
		{
			name:      "rate 35%: income=10,000,000 deduction.personal=0; expect tax=3,110,000",
			info:      TaxInformation{TotalIncome: 10_000_000.0 + defaultDeduction.Personal},
			deduction: defaultDeduction,
			want:      TaxResult{Tax: 3_110_000.0},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := CalculateTax(tc.info, tc.deduction)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCalculateTax_Error_InputOnlyTotalIncome(t *testing.T) {
	// Arrange
	invalidDeduction := Deduction{
		Personal: ConstraintMaxPersonalDeduction + 0.1,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	got, err := CalculateTax(TaxInformation{TotalIncome: 100_000.0}, invalidDeduction)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidDeduction)
	assert.Equal(t, TaxResult{}, got)
}
