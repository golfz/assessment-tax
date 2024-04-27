//go:build unit

package deduction

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate_Success(t *testing.T) {
	defaultDeduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}
	// Arrange
	testCases := []struct {
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
				Personal: MinPersonalDeduction + 0.1,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
		},
		{
			name: "personal deduction = max",
			deduction: Deduction{
				Personal: MaxPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
		},
		// KReceipt deduction
		{
			name: "KReceipt deduction = min",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: MinKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
		},
		{
			name: "KReceipt deduction = max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: MaxKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
		},
		// Donation deduction
		{
			name: "Donation deduction = max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: defaultDeduction.KReceipt,
				Donation: MaxDonationDeduction,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotError := tc.deduction.Validate()

			// Assert
			assert.NoError(t, gotError)
		})
	}
}

func TestValidate_Error(t *testing.T) {
	defaultDeduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}
	// Arrange
	testCases := []struct {
		name       string
		deduction  Deduction
		wantErrors []error
	}{
		// personal deduction
		{
			name: "personal deduction < min",
			deduction: Deduction{
				Personal: MinPersonalDeduction,
				KReceipt: defaultDeduction.KReceipt,
				Donation: defaultDeduction.Donation,
			},
			wantErrors: []error{ErrInvalidPersonalDeduction},
		},
		{
			name: "personal deduction > max",
			deduction: Deduction{
				Personal: MaxPersonalDeduction + 0.1,
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
				KReceipt: MinKReceiptDeduction,
				Donation: defaultDeduction.Donation,
			},
			wantErrors: []error{ErrInvalidKReceiptDeduction},
		},
		{
			name: "KReceipt deduction > max",
			deduction: Deduction{
				Personal: defaultDeduction.Personal,
				KReceipt: MaxKReceiptDeduction + 0.1,
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
				KReceipt: MaxKReceiptDeduction + 0.1,
				Donation: MaxDonationDeduction + 0.1,
			},
			wantErrors: []error{ErrInvalidDonationDeduction},
		},
		// Multiple errors
		{
			name: "personal deduction > max, KReceipt deduction > max",
			deduction: Deduction{
				Personal: MaxPersonalDeduction + 0.1,
				KReceipt: MaxKReceiptDeduction + 0.1,
				Donation: defaultDeduction.Donation,
			},
			wantErrors: []error{ErrInvalidPersonalDeduction, ErrInvalidKReceiptDeduction},
		},
		{
			name: "personal deduction > max, KReceipt deduction > max, Donation deduction > max",
			deduction: Deduction{
				Personal: MaxPersonalDeduction + 0.1,
				KReceipt: MaxKReceiptDeduction + 0.1,
				Donation: MaxDonationDeduction + 0.1,
			},
			wantErrors: []error{ErrInvalidPersonalDeduction, ErrInvalidKReceiptDeduction, ErrInvalidDonationDeduction},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotError := tc.deduction.Validate()

			// Assert
			assert.Error(t, gotError)
			for _, wantError := range tc.wantErrors {
				assert.ErrorIs(t, gotError, wantError)
			}
		})
	}
}
