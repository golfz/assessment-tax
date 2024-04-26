//go:build integration

package tax_test

import (
	"encoding/json"
	"github.com/golfz/assessment-tax/config"
	"github.com/golfz/assessment-tax/postgres"
	"github.com/golfz/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCalculateTaxIntegration_Success(t *testing.T) {
	testcases := []struct {
		name          string
		taxInfo       tax.TaxInformation
		wantTaxResult tax.TaxResult
	}{
		{
			name: "EXP01: basic income, no WHT, no Allowance; expect tax",
			taxInfo: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 29_000.0, TaxRefund: 0.0},
		},
		{
			name: "EXP02: Income and WHT, no Allowance; expect tax",
			taxInfo: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         25_000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 4_000.0, TaxRefund: 0.0},
		},
		{
			name: "EXP03: Income and Allowance, no WHT; expect tax",
			taxInfo: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 19_000.0, TaxRefund: 0.0},
		},
		{
			name: "One Allowance, tax payable > WHT; expect tax",
			taxInfo: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         15_000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 4_000.0, TaxRefund: 0.0},
		},
		{
			name: "One Allowance, tax payable = WHT; expect tax=0",
			taxInfo: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         19_000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 0.0, TaxRefund: 0.0},
		},
		{
			name: "One Allowance, tax payable < WHT; expect taxRefund",
			taxInfo: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         29_000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
		},
		{
			name: "Multi Allowance, tax payable > WHT; expect tax",
			taxInfo: tax.TaxInformation{
				TotalIncome: 600_000.0,
				WHT:         15_000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeKReceipt, Amount: 40_000.0},
					{Type: tax.AllowanceTypeKReceipt, Amount: 30_000.0},
					{Type: tax.AllowanceTypeDonation, Amount: 80_000.0},
					{Type: tax.AllowanceTypeDonation, Amount: 70_000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 9_000.0, TaxRefund: 0.0},
		},
		{
			name: "Multi Allowance, tax payable = WHT; expect tax=0",
			taxInfo: tax.TaxInformation{
				TotalIncome: 600_000.0,
				WHT:         24_000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeKReceipt, Amount: 40_000.0},
					{Type: tax.AllowanceTypeKReceipt, Amount: 30_000.0},
					{Type: tax.AllowanceTypeDonation, Amount: 80_000.0},
					{Type: tax.AllowanceTypeDonation, Amount: 70_000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 0.0, TaxRefund: 0.0},
		},
		{
			name: "Multi Allowance, tax payable < WHT; expect taxRefund>0",
			taxInfo: tax.TaxInformation{
				TotalIncome: 600_000.0,
				WHT:         34_000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeKReceipt, Amount: 40_000.0},
					{Type: tax.AllowanceTypeKReceipt, Amount: 30_000.0},
					{Type: tax.AllowanceTypeDonation, Amount: 80_000.0},
					{Type: tax.AllowanceTypeDonation, Amount: 70_000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			cfg := config.NewWith(os.Getenv)

			pg, err := postgres.New(cfg.DatabaseURL)
			if err != nil {
				t.Errorf("failed to connect to database: %v", err)
			}

			e := echo.New()
			hTax := tax.New(pg)
			e.POST("/tax/calculations", hTax.CalculateTaxHandler)

			var bReader io.Reader
			b, _ := json.Marshal(tc.taxInfo)
			bReader = strings.NewReader(string(b))
			req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bReader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// Act
			e.ServeHTTP(rec, req)

			// Assert
			assert.Equal(t, http.StatusOK, rec.Code)

			var gotTaxResult tax.TaxResult
			if err := json.Unmarshal(rec.Body.Bytes(), &gotTaxResult); err != nil {
				t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
			}
			assert.Equal(t, tc.wantTaxResult.Tax, gotTaxResult.Tax)
			assert.Equal(t, tc.wantTaxResult.TaxRefund, gotTaxResult.TaxRefund)
		})
	}
}

func TestCalculateTaxIntegration_WithTaxLevel_Success(t *testing.T) {
	// Arrange
	testcases := []struct {
		name          string
		taxInfo       tax.TaxInformation
		wantTaxResult tax.TaxResult
		wantTaxLevels []float64
	}{
		{
			name: "EXP04: net-income=340,000 (rate=10%); expect tax=19,000",
			taxInfo: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 200000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 19_000.0, TaxRefund: 0.0},
			wantTaxLevels: []float64{0.0, 19_000.0, 0.0, 0.0, 0.0},
		},
		{
			name: "net-income=100,000 (rate=0%); expect tax=0",
			taxInfo: tax.TaxInformation{
				TotalIncome: 260_000.0,
				WHT:         0.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 200000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 0.0, TaxRefund: 0.0},
			wantTaxLevels: []float64{0.0, 0.0, 0.0, 0.0, 0.0},
		},
		{
			name: "net-income=3,000,000 (rate=35%); expect tax=660,000",
			taxInfo: tax.TaxInformation{
				TotalIncome: 3_160_000.0,
				WHT:         0.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 200000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 660_000.0, TaxRefund: 0.0},
			wantTaxLevels: []float64{0.0, 35_000.0, 75_000.0, 200_000.0, 350_000.0},
		},
		{
			name: "net-income=3,000,000 (rate=35%) wht=700,000; expect taxRefund=40,000",
			taxInfo: tax.TaxInformation{
				TotalIncome: 3_160_000.0,
				WHT:         700_000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 200000.0},
				},
			},
			wantTaxResult: tax.TaxResult{Tax: 0.0, TaxRefund: 40_000.0},
			wantTaxLevels: []float64{0.0, 35_000.0, 75_000.0, 200_000.0, 350_000.0},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			cfg := config.NewWith(os.Getenv)

			pg, err := postgres.New(cfg.DatabaseURL)
			if err != nil {
				t.Errorf("failed to connect to database: %v", err)
			}

			e := echo.New()
			hTax := tax.New(pg)
			e.POST("/tax/calculations", hTax.CalculateTaxHandler)

			var bReader io.Reader
			b, _ := json.Marshal(tc.taxInfo)
			bReader = strings.NewReader(string(b))
			req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bReader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// Act
			e.ServeHTTP(rec, req)

			// Assert
			assert.Equal(t, http.StatusOK, rec.Code)

			var gotTaxResult tax.TaxResult
			if err := json.Unmarshal(rec.Body.Bytes(), &gotTaxResult); err != nil {
				t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
			}
			assert.Equal(t, tc.wantTaxResult.Tax, gotTaxResult.Tax)
			assert.Equal(t, tc.wantTaxResult.TaxRefund, gotTaxResult.TaxRefund)

			for i, wantTax := range tc.wantTaxLevels {
				assert.Equal(t, wantTax, gotTaxResult.TaxLevels[i].Tax)
			}
		})
	}
}
