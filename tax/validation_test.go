//go:build unit

package tax

import (
	"github.com/golfz/assessment-tax/rule"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateTaxInformation_Success(t *testing.T) {
	// Arrange
	testcases := []struct {
		name    string
		taxInfo TaxInformation
	}{
		{
			name:    "total income = 0",
			taxInfo: TaxInformation{TotalIncome: 0.0},
		},
		{
			name:    "total income = 100,000",
			taxInfo: TaxInformation{TotalIncome: 100_000.0},
		},
		{
			name:    "WHT = 0",
			taxInfo: TaxInformation{TotalIncome: 100_000.0, WHT: 0.0},
		},
		{
			name:    "WHT < total income",
			taxInfo: TaxInformation{TotalIncome: 100_000.0, WHT: 10_000.0},
		},
		{
			name:    "WHT = total income",
			taxInfo: TaxInformation{TotalIncome: 100_000.0, WHT: 100_000.0},
		},
		{
			name: "allowance = 0 or allowance > 0",
			taxInfo: TaxInformation{
				TotalIncome: 100_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
					{Type: AllowanceTypeKReceipt, Amount: 10_000.0},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotError := validateTaxInformation(tc.taxInfo)

			// Assert
			assert.NoError(t, gotError)
		})
	}
}

func TestValidateTaxInformation_Error(t *testing.T) {
	// Arrange
	testcases := []struct {
		name       string
		taxInfo    TaxInformation
		wantErrors []error
	}{
		{
			name:       "total income < 0",
			taxInfo:    TaxInformation{TotalIncome: -1.0},
			wantErrors: []error{ErrInvalidTotalIncome},
		},
		{
			name:       "WHT < 0",
			taxInfo:    TaxInformation{TotalIncome: 100_000.0, WHT: -1.0},
			wantErrors: []error{ErrInvalidWHT},
		},
		{
			name:       "WHT > income",
			taxInfo:    TaxInformation{TotalIncome: 100_000.0, WHT: 200_000.0},
			wantErrors: []error{ErrInvalidWHT},
		},
		{
			name:       "total income < 0 and WHT < 0",
			taxInfo:    TaxInformation{TotalIncome: -1.0, WHT: -1.0},
			wantErrors: []error{ErrInvalidTotalIncome, ErrInvalidWHT},
		},
		{
			name: "some allowance < 0",
			taxInfo: TaxInformation{
				TotalIncome: 100_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: -1.0},
					{Type: AllowanceTypeDonation, Amount: 0.0},
					{Type: AllowanceTypeKReceipt, Amount: 10_000.0},
				},
			},
			wantErrors: []error{ErrInvalidAllowanceAmount},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotError := validateTaxInformation(tc.taxInfo)

			// Assert
			assert.Error(t, gotError)

			for _, wantError := range tc.wantErrors {
				assert.ErrorIs(t, gotError, wantError)
			}
		})
	}
}

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
				Personal: rule.MinPersonalDeduction + 0.1,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
		},
		{
			name: "personal deduction = max",
			deduction: Deduction{
				Personal: rule.MaxPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
		},
		// KReceipt deduction
		{
			name: "KReceipt deduction = min",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: rule.MinKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
		},
		{
			name: "KReceipt deduction = max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: rule.MaxKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
		},
		// Donation deduction
		{
			name: "Donation deduction = max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: rule.MaxDonationDeduction,
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
		name       string
		deduction  Deduction
		wantErrors []error
	}{
		// personal deduction
		{
			name: "personal deduction < min",
			deduction: Deduction{
				Personal: rule.MinPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantErrors: []error{ErrInvalidPersonalDeduction},
		},
		{
			name: "personal deduction > max",
			deduction: Deduction{
				Personal: rule.MaxPersonalDeduction + 0.1,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantErrors: []error{ErrInvalidPersonalDeduction},
		},
		// KReceipt deduction
		{
			name: "KReceipt deduction < min",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: rule.MinKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
			wantErrors: []error{ErrInvalidKReceiptDeduction},
		},
		{
			name: "KReceipt deduction > max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: rule.MaxKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
			wantErrors: []error{ErrInvalidKReceiptDeduction},
		},
		// Donation deduction
		{
			name: "Donation deduction > max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				//KReceipt: defaultDeduction.KReceipt,
				KReceipt: rule.MaxKReceiptDeduction + 0.1,
				Donation: rule.MaxDonationDeduction + 0.1,
			},
			wantErrors: []error{ErrInvalidDonationDeduction},
		},
		// Multiple errors
		{
			name: "personal deduction > max, KReceipt deduction > max",
			deduction: Deduction{
				Personal: rule.MaxPersonalDeduction + 0.1,
				KReceipt: rule.MaxKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
			wantErrors: []error{ErrInvalidPersonalDeduction, ErrInvalidKReceiptDeduction},
		},
		{
			name: "personal deduction > max, KReceipt deduction > max, Donation deduction > max",
			deduction: Deduction{
				Personal: rule.MaxPersonalDeduction + 0.1,
				KReceipt: rule.MaxKReceiptDeduction + 0.1,
				Donation: rule.MaxDonationDeduction + 0.1,
			},
			wantErrors: []error{ErrInvalidPersonalDeduction, ErrInvalidKReceiptDeduction, ErrInvalidDonationDeduction},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotError := validateDeduction(tc.deduction)

			// Assert
			assert.Error(t, gotError)
			for _, wantError := range tc.wantErrors {
				assert.ErrorIs(t, gotError, wantError)
			}
		})
	}
}
