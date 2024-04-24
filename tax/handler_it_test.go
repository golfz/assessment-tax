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

func TestCalculateTax_Success_Integration(t *testing.T) {
	t.Run("CalculateTax_Success_Integration", func(t *testing.T) {
		// Arrange
		cfg := config.NewWith(os.Getenv)

		pg, err := postgres.New(cfg.DatabaseURL)
		if err != nil {
			t.Errorf("failed to connect to database: %v", err)
		}

		e := echo.New()
		hTax := tax.New(pg)
		e.POST("/tax/calculations", hTax.CalculateTaxHandler)

		reqInfoTax := tax.TaxInformation{
			TotalIncome: 500_000.0,
			WHT:         0.0,
			Allowances: []tax.Allowance{
				{Type: tax.AllowanceTypeDonation, Amount: 0.0},
			},
		}
		var bReader io.Reader
		b, _ := json.Marshal(reqInfoTax)
		bReader = strings.NewReader(string(b))
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bReader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		wantTax := 29_000.0
		wantCode := http.StatusOK

		// Act
		e.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, wantCode, rec.Code)

		var got tax.TaxResult
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
		}
		assert.Equal(t, wantTax, got.Tax)
	})

}
