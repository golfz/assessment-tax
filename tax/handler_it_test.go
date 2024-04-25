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

func TestCalculateTaxIntegration_Success_200_OK_With_Tax_And_TaxRefund(t *testing.T) {
	testcases := []struct {
		name          string
		info          tax.TaxInformation
		wantTax       float64
		wantTaxRefund float64
	}{
		{
			name: "EXP01: Only TotalIncome; expect tax=29_000.0, taxRefund=0.0",
			info: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTax:       29_000.0,
			wantTaxRefund: 0.0,
		},
		{
			name: "EXP02: Income and WHT; expect tax=4_000.0, taxRefund=0.0",
			info: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         25_000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTax:       4_000.0,
			wantTaxRefund: 0.0,
		},
		{
			name: "EXP03: Income and Allowance; expect tax=19_000.0, taxRefund=0.0",
			info: tax.TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []tax.Allowance{
					{
						Type:   tax.AllowanceTypeDonation,
						Amount: 200_000.0,
					},
				},
			},
			wantTax:       19_000.0,
			wantTaxRefund: 0.0,
		},
		{
			name:          "tax > WHT; expect tax>0",
			info:          tax.TaxInformation{TotalIncome: 500_000.0, WHT: 25_000.0},
			wantTax:       4_000.0,
			wantTaxRefund: 0.0,
		},
		{
			name:          "tax = WHT; expect tax=0",
			info:          tax.TaxInformation{TotalIncome: 500_000.0, WHT: 29_000.0},
			wantTax:       0.0,
			wantTaxRefund: 0.0,
		},
		{
			name:          "tax < WHT; expect taxRefund>0",
			info:          tax.TaxInformation{TotalIncome: 500_000.0, WHT: 39_000.0},
			wantTax:       0.0,
			wantTaxRefund: 10_000.0,
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
			b, _ := json.Marshal(tc.info)
			bReader = strings.NewReader(string(b))
			req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bReader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// Act
			e.ServeHTTP(rec, req)

			// Assert
			assert.Equal(t, http.StatusOK, rec.Code)

			var got tax.TaxResult
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
			}
			assert.Equal(t, tc.wantTax, got.Tax)
			assert.Equal(t, tc.wantTaxRefund, got.TaxRefund)
		})
	}

}
