//go:build unit

package tax

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestCalculateTax_Success_from_TotalIncome_and_WHT(t *testing.T) {
	// Arrange
	deduction := Deduction{
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
			name:      "tax > WHT; expect tax>0",
			info:      TaxInformation{TotalIncome: 500_000.0, WHT: 25_000.0},
			deduction: deduction,
			want:      TaxResult{Tax: 4000.0, TaxRefund: 0.0},
		},
		{
			name:      "tax = WHT; expect tax=0",
			info:      TaxInformation{TotalIncome: 500_000.0, WHT: 29_000.0},
			deduction: deduction,
			want:      TaxResult{Tax: 0.0, TaxRefund: 0.0},
		},
		{
			name:      "tax < WHT; expect taxRefund>0",
			info:      TaxInformation{TotalIncome: 500_000.0, WHT: 39_000.0},
			deduction: deduction,
			want:      TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
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

func TestCalculateTax_Error_InvalidTaxInformation(t *testing.T) {
	deduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	t.Run("total income < 0", func(t *testing.T) {
		// Arrange
		invalidInfo := TaxInformation{TotalIncome: -1.0}

		// Act
		got, err := CalculateTax(invalidInfo, deduction)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidTaxInformation)
		assert.ErrorIs(t, err, ErrInvalidTotalIncome)
		assert.NotErrorIs(t, err, ErrInvalidWHT)
		assert.NotErrorIs(t, err, ErrInvalidAllowanceAmount)
		assert.Equal(t, TaxResult{}, got)
	})

	t.Run("WHT < 0", func(t *testing.T) {
		// Arrange
		invalidInfo := TaxInformation{TotalIncome: 100_000.0, WHT: -1.0}

		// Act
		got, err := CalculateTax(invalidInfo, deduction)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidTaxInformation)
		assert.ErrorIs(t, err, ErrInvalidWHT)
		assert.NotErrorIs(t, err, ErrInvalidTotalIncome)
		assert.NotErrorIs(t, err, ErrInvalidAllowanceAmount)
		assert.Equal(t, TaxResult{}, got)
	})

	t.Run("WHT > total income", func(t *testing.T) {
		// Arrange
		invalidInfo := TaxInformation{TotalIncome: 100_000.0, WHT: 200_000.0}

		// Act
		got, err := CalculateTax(invalidInfo, Deduction{})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidTaxInformation)
		assert.ErrorIs(t, err, ErrInvalidWHT)
		assert.NotErrorIs(t, err, ErrInvalidTotalIncome)
		assert.NotErrorIs(t, err, ErrInvalidAllowanceAmount)
		assert.Equal(t, TaxResult{}, got)
	})

	t.Run("some allowance < 0", func(t *testing.T) {
		// Arrange
		invalidInfo := TaxInformation{
			TotalIncome: 100_000.0,
			Allowances: []Allowance{
				{Type: AllowanceTypeDonation, Amount: -1.0},
				{Type: AllowanceTypeDonation, Amount: 0.0},
				{Type: AllowanceTypeKReceipt, Amount: 10_000.0},
			},
		}

		// Act
		got, err := CalculateTax(invalidInfo, deduction)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidTaxInformation)
		assert.ErrorIs(t, err, ErrInvalidAllowanceAmount)
		assert.NotErrorIs(t, err, ErrInvalidTotalIncome)
		assert.NotErrorIs(t, err, ErrInvalidWHT)
		assert.Equal(t, TaxResult{}, got)
	})
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
