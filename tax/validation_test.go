//go:build unit

package tax

import (
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
		name    string
		taxInfo TaxInformation
		error   []error
	}{
		{
			name:    "total income < 0",
			taxInfo: TaxInformation{TotalIncome: -1.0},
			error:   []error{ErrInvalidTotalIncome},
		},
		{
			name:    "WHT < 0",
			taxInfo: TaxInformation{TotalIncome: 100_000.0, WHT: -1.0},
			error:   []error{ErrInvalidWHT},
		},
		{
			name:    "WHT > income",
			taxInfo: TaxInformation{TotalIncome: 100_000.0, WHT: 200_000.0},
			error:   []error{ErrInvalidWHT},
		},
		{
			name:    "total income < 0 and WHT < 0",
			taxInfo: TaxInformation{TotalIncome: -1.0, WHT: -1.0},
			error:   []error{ErrInvalidTotalIncome, ErrInvalidWHT},
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
			error: []error{ErrInvalidAllowanceAmount},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotError := validateTaxInformation(tc.taxInfo)

			// Assert
			assert.Error(t, gotError)

			for _, wantErr := range tc.error {
				assert.ErrorIs(t, gotError, wantErr)
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
