//go:build integration

package admin_test

import (
	"encoding/json"
	"github.com/golfz/assessment-tax/admin"
	"github.com/golfz/assessment-tax/config"
	"github.com/golfz/assessment-tax/deduction"
	"github.com/golfz/assessment-tax/postgres"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setup(t *testing.T) func() {
	return func() {
		// Arrange
		cfg := config.NewWith(os.Getenv)

		pg, err := postgres.New(cfg.DatabaseURL)
		if err != nil {
			t.Errorf("failed to connect to database: %v", err)
		}

		e := echo.New()
		hAdmin := admin.New(pg)
		e.POST("/admin/deductions/personal", hAdmin.SetPersonalDeductionHandler)

		input := admin.Input{
			Amount: deduction.DefaultPersonalDeduction,
		}
		var bReader io.Reader
		b, _ := json.Marshal(input)
		bReader = strings.NewReader(string(b))
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bReader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Act
		e.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestSetPersonalDeductionIntegration_Success(t *testing.T) {
	testCases := []struct {
		name  string
		input admin.Input
		want  admin.PersonalDeduction
	}{
		{
			name: "setting with default personal deduction",
			input: admin.Input{
				Amount: deduction.DefaultPersonalDeduction,
			},
			want: admin.PersonalDeduction{
				PersonalDeduction: deduction.DefaultPersonalDeduction,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			teardown := setup(t)
			defer teardown()

			// Arrange
			cfg := config.NewWith(os.Getenv)

			pg, err := postgres.New(cfg.DatabaseURL)
			if err != nil {
				t.Errorf("failed to connect to database: %v", err)
			}

			e := echo.New()
			hAdmin := admin.New(pg)
			e.POST("/admin/deductions/personal", hAdmin.SetPersonalDeductionHandler)

			var bReader io.Reader
			b, _ := json.Marshal(tc.input)
			bReader = strings.NewReader(string(b))
			req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bReader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// Act
			e.ServeHTTP(rec, req)

			// Assert
			assert.Equal(t, http.StatusOK, rec.Code)

			var got admin.PersonalDeduction
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
