//go:build unit

package tax

import (
	"github.com/golfz/assessment-tax/rule"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateNetIncome(t *testing.T) {
	// Arrange
	testcases := []struct {
		name              string
		totalIncome       float64
		personalDeduction float64
		totalAllowance    float64
		want              float64
	}{
		{
			name:              "income=100,000 personal=60,000 allowance=0; expect net=40,000",
			totalIncome:       100_000.0,
			personalDeduction: 60_000.0,
			totalAllowance:    0.0,
			want:              40_000.0,
		},
		{
			name:              "income=100,000 personal=60,000 allowance=10,000; expect net=30,000",
			totalIncome:       100_000.0,
			personalDeduction: 60_000.0,
			totalAllowance:    10_000.0,
			want:              30_000.0,
		},
		{
			name:              "income=100,000 personal=60,000 allowance=40,000; expect net=0",
			totalIncome:       100_000.0,
			personalDeduction: 60_000.0,
			totalAllowance:    40_000.0,
			want:              0.0,
		},
		{
			name:              "income=100,000 personal=60,000 allowance=100,000; expect net=0",
			totalIncome:       100_000.0,
			personalDeduction: 60_000.0,
			totalAllowance:    100_000.0,
			want:              0.0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := calculateNetIncome(tc.totalIncome, tc.personalDeduction, tc.totalAllowance)

			// Assert
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCalculateTax_ByRateFromIncomeOnly_ExpectSuccess(t *testing.T) {
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
			assert.Equal(t, tc.want.Tax, got.Tax)
			assert.Equal(t, tc.want.TaxRefund, got.TaxRefund)
		})
	}
}

func TestCalculateTax_Success(t *testing.T) {
	// Arrange
	deduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}
	testcases := []struct {
		name    string
		taxInfo TaxInformation
		want    TaxResult
	}{
		{
			name: "EXP01: only income, net-income=290,000 (rate=10%); expect tax=29,000",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			want: TaxResult{Tax: 29_000.0},
		},
		{
			name: "only income, net-income=150,000 (rate=0%); expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 210_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			want: TaxResult{Tax: 0.0},
		},
		{
			name: "only income, net-income=0 (rate=0%); expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 60_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			want: TaxResult{Tax: 0.0},
		},
		{
			name: "EXP02: tax-payable>wht; expect tax=4,000",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         25_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			want: TaxResult{Tax: 4_000.0},
		},
		{
			name: "tax-payable=wht; expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         29_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			want: TaxResult{Tax: 0.0},
		},
		{
			name: "tax-payable<wht; expect taxRefund=10,000",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         39_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			want: TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
		},
		{
			name: "EXP03: income=500,000 donation=200,000; expect tax=19,000",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			want: TaxResult{Tax: 19_000.0},
		},
		{
			name: "income=500,000 wht=tax-payable donation=200,000; expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         19_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			want: TaxResult{Tax: 0.0},
		},
		{
			name: "income=500,000 wht>tax-payable donation=200,000; expect taxRefund=10,000",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         29_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			want: TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
		},
		{
			name: "netIncome=0: income=200,000 deduction.personal=60,000 allowance=140,000; expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 200_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 100_000.0},
					{Type: AllowanceTypeKReceipt, Amount: 40_000.0},
				},
			},
			want: TaxResult{Tax: 0.0},
		},
		{
			name: "netIncome=0: income=200,000 wht=10,000 deduction.personal=60,000 allowance=140,000; expect taxRefund=10,000",
			taxInfo: TaxInformation{
				TotalIncome: 200_000.0,
				WHT:         10_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 100_000.0},
					{Type: AllowanceTypeKReceipt, Amount: 40_000.0},
				},
			},
			want: TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
		},
		{
			name: "netIncome<0: income=150,000 deduction.personal=60,000 allowance=140,000; expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 150_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 100_000.0},
					{Type: AllowanceTypeKReceipt, Amount: 40_000.0},
				},
			},
			want: TaxResult{Tax: 0.0},
		},
		{
			name: "netIncome<0: income=150,000 wht=10,000 deduction.personal=60,000 allowance=140,000; expect taxRefund=10,000",
			taxInfo: TaxInformation{
				TotalIncome: 150_000.0,
				WHT:         10_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 100_000.0},
					{Type: AllowanceTypeKReceipt, Amount: 40_000.0},
				},
			},
			want: TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
		},
		{
			name: "Multi Allowance, tax payable > WHT; expect tax",
			taxInfo: TaxInformation{
				TotalIncome: 600_000.0,
				WHT:         15_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeKReceipt, Amount: 40_000.0},
					{Type: AllowanceTypeKReceipt, Amount: 30_000.0},
					{Type: AllowanceTypeDonation, Amount: 80_000.0},
					{Type: AllowanceTypeDonation, Amount: 70_000.0},
				},
			},
			want: TaxResult{Tax: 9_000.0, TaxRefund: 0.0},
		},
		{
			name: "Multi Allowance, tax payable = WHT; expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 600_000.0,
				WHT:         24_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeKReceipt, Amount: 40_000.0},
					{Type: AllowanceTypeKReceipt, Amount: 30_000.0},
					{Type: AllowanceTypeDonation, Amount: 80_000.0},
					{Type: AllowanceTypeDonation, Amount: 70_000.0},
				},
			},
			want: TaxResult{Tax: 0.0, TaxRefund: 0.0},
		},
		{
			name: "Multi Allowance, tax payable < WHT; expect taxRefund>0",
			taxInfo: TaxInformation{
				TotalIncome: 600_000.0,
				WHT:         34_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeKReceipt, Amount: 40_000.0},
					{Type: AllowanceTypeKReceipt, Amount: 30_000.0},
					{Type: AllowanceTypeDonation, Amount: 80_000.0},
					{Type: AllowanceTypeDonation, Amount: 70_000.0},
				},
			},
			want: TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := CalculateTax(tc.taxInfo, deduction)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.want.Tax, got.Tax)
			assert.Equal(t, tc.want.TaxRefund, got.TaxRefund)
		})
	}
}

func TestCalculateTax_WithTaxLevel(t *testing.T) {
	// Arrange
	deduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}
	testcases := []struct {
		name          string
		taxInfo       TaxInformation
		wantTaxResult TaxResult
		wantTaxLevels []float64
	}{
		{
			name: "EXP04: net-income=340,000 (rate=10%); expect tax=19,000",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200000.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 19_000.0, TaxRefund: 0.0},
			wantTaxLevels: []float64{0.0, 19_000.0, 0.0, 0.0, 0.0},
		},
		{
			name: "net-income=100,000 (rate=0%); expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 260_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200000.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 0.0, TaxRefund: 0.0},
			wantTaxLevels: []float64{0.0, 0.0, 0.0, 0.0, 0.0},
		},
		{
			name: "net-income=3,000,000 (rate=35%); expect tax=660,000",
			taxInfo: TaxInformation{
				TotalIncome: 3_160_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200000.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 660_000.0, TaxRefund: 0.0},
			wantTaxLevels: []float64{0.0, 35_000.0, 75_000.0, 200_000.0, 350_000.0},
		},
		{
			name: "net-income=3,000,000 (rate=35%) wht=700,000; expect taxRefund=40,000",
			taxInfo: TaxInformation{
				TotalIncome: 3_160_000.0,
				WHT:         700_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200000.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 0.0, TaxRefund: 40_000.0},
			wantTaxLevels: []float64{0.0, 35_000.0, 75_000.0, 200_000.0, 350_000.0},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := CalculateTax(tc.taxInfo, deduction)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.wantTaxResult.Tax, got.Tax)
			assert.Equal(t, tc.wantTaxResult.TaxRefund, got.TaxRefund)
			for i, wantTax := range tc.wantTaxLevels {
				assert.Equal(t, wantTax, got.TaxLevels[i].Tax)
			}
		})
	}
}

func TestCalculateTax_FromInvalidTaxInformation_Error(t *testing.T) {
	// Arrange
	deduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	testcases := []struct {
		name           string
		taxInformation TaxInformation
		wantErrors     []error
		unwantedErrors []error
	}{
		{
			name:           "total income < 0",
			taxInformation: TaxInformation{TotalIncome: -1.0},
			wantErrors:     []error{ErrInvalidTotalIncome, ErrInvalidTaxInformation},
			unwantedErrors: []error{ErrInvalidWHT, ErrInvalidAllowanceAmount},
		},
		{
			name:           "WHT < 0",
			taxInformation: TaxInformation{TotalIncome: 100_000.0, WHT: -1.0},
			wantErrors:     []error{ErrInvalidWHT, ErrInvalidTaxInformation},
			unwantedErrors: []error{ErrInvalidTotalIncome, ErrInvalidAllowanceAmount},
		},
		{
			name:           "WHT > income",
			taxInformation: TaxInformation{TotalIncome: 100_000.0, WHT: 200_000.0},
			wantErrors:     []error{ErrInvalidWHT, ErrInvalidTaxInformation},
			unwantedErrors: []error{ErrInvalidTotalIncome, ErrInvalidAllowanceAmount},
		},
		{
			name:           "total income < 0 and WHT < 0",
			taxInformation: TaxInformation{TotalIncome: -1.0, WHT: -1.0},
			wantErrors:     []error{ErrInvalidTotalIncome, ErrInvalidWHT, ErrInvalidTaxInformation},
			unwantedErrors: []error{ErrInvalidAllowanceAmount},
		},
		{
			name: "some allowance < 0",
			taxInformation: TaxInformation{
				TotalIncome: 100_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: -1.0},
					{Type: AllowanceTypeDonation, Amount: 0.0},
					{Type: AllowanceTypeKReceipt, Amount: 10_000.0},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := CalculateTax(tc.taxInformation, deduction)

			// Assert
			assert.Error(t, err)
			for _, wantError := range tc.wantErrors {
				assert.ErrorIs(t, err, wantError)
			}
			for _, unwantedError := range tc.unwantedErrors {
				assert.NotErrorIs(t, err, unwantedError)
			}
			assert.Equal(t, TaxResult{}, got)
		})
	}
}

func TestCalculateTax_FromInvalidDeduction_Error(t *testing.T) {
	t.Run("personal deduction > max", func(t *testing.T) {
		// Arrange
		invalidDeduction := Deduction{
			Personal: rule.MaxPersonalDeduction + 0.1,
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
			KReceipt: rule.MaxKReceiptDeduction + 0.1,
			Donation: rule.MaxDonationDeduction + 0.1,
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
