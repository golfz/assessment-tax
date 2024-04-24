//go:build unit

package tax

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateDeduction_Success(t *testing.T) {
	defaultDeduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}
	// Arrange
	testcases := []struct {
		name      string
		deduction Deduction
	}{
		// Default deduction
		{
			name: "default deduction",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
		},
		// personal deduction
		{
			name: "personal deduction = min",
			deduction: Deduction{
				Personal: ConstraintMinPersonalDeduction + 0.1,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
		},
		{
			name: "personal deduction = max",
			deduction: Deduction{
				Personal: ConstraintMaxPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
		},
		// KReceipt deduction
		{
			name: "KReceipt deduction = min",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMinKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
		},
		{
			name: "KReceipt deduction = max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMaxKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
		},
		// Donation deduction
		{
			name: "Donation deduction = max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: ConstraintMaxDonationDeduction,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotError := validateDeduction(tc.deduction)

			// Assert
			assert.NoError(t, gotError)
		})
	}
}

func TestValidateDeduction_Error(t *testing.T) {
	defaultDeduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}
	// Arrange
	testcases := []struct {
		name      string
		deduction Deduction
		wantError []error
	}{
		// personal deduction
		{
			name: "personal deduction < min",
			deduction: Deduction{
				Personal: ConstraintMinPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantError: []error{ErrInvalidPersonalDeduction},
		},
		{
			name: "personal deduction > max",
			deduction: Deduction{
				Personal: ConstraintMaxPersonalDeduction + 0.1,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantError: []error{ErrInvalidPersonalDeduction},
		},
		// KReceipt deduction
		{
			name: "KReceipt deduction < min",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMinKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
			wantError: []error{ErrInvalidKReceiptDeduction},
		},
		{
			name: "KReceipt deduction > max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: ConstraintMaxKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
			wantError: []error{ErrInvalidKReceiptDeduction},
		},
		// Donation deduction
		{
			name: "Donation deduction > max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				//KReceipt: defaultDeduction.KReceipt,
				KReceipt: ConstraintMaxKReceiptDeduction + 0.1,
				Donation: ConstraintMaxDonationDeduction + 0.1,
			},
			wantError: []error{ErrInvalidDonationDeduction},
		},
		// Multiple errors
		{
			name: "personal deduction > max, KReceipt deduction > max",
			deduction: Deduction{
				Personal: ConstraintMaxPersonalDeduction + 0.1,
				KReceipt: ConstraintMaxKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
			wantError: []error{ErrInvalidPersonalDeduction, ErrInvalidKReceiptDeduction},
		},
		{
			name: "personal deduction > max, KReceipt deduction > max, Donation deduction > max",
			deduction: Deduction{
				Personal: ConstraintMaxPersonalDeduction + 0.1,
				KReceipt: ConstraintMaxKReceiptDeduction + 0.1,
				Donation: ConstraintMaxDonationDeduction + 0.1,
			},
			wantError: []error{ErrInvalidPersonalDeduction, ErrInvalidKReceiptDeduction, ErrInvalidDonationDeduction},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotError := validateDeduction(tc.deduction)

			// Assert
			assert.Error(t, gotError)
			for _, wantErr := range tc.wantError {
				assert.ErrorIs(t, gotError, wantErr)
			}
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

func TestCalculateTax_Error_InvalidDeduction(t *testing.T) {
	t.Run("personal deduction > max", func(t *testing.T) {
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
		assert.ErrorIs(t, err, ErrInvalidPersonalDeduction)
		assert.NotErrorIs(t, err, ErrInvalidKReceiptDeduction)
		assert.NotErrorIs(t, err, ErrInvalidDonationDeduction)
		assert.Equal(t, TaxResult{}, got)
	})

	t.Run("KReceipt deduction > max and donation deduction > max", func(t *testing.T) {
		// Arrange
		invalidDeduction := Deduction{
			Personal: 60_000.0,
			KReceipt: ConstraintMaxKReceiptDeduction + 0.1,
			Donation: ConstraintMaxDonationDeduction + 0.1,
		}

		// Act
		got, err := CalculateTax(TaxInformation{TotalIncome: 100_000.0}, invalidDeduction)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidDeduction)
		assert.ErrorIs(t, err, ErrInvalidKReceiptDeduction)
		assert.ErrorIs(t, err, ErrInvalidDonationDeduction)
		assert.NotErrorIs(t, err, ErrInvalidPersonalDeduction)
		assert.Equal(t, TaxResult{}, got)
	})
}
