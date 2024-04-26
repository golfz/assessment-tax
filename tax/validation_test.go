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
