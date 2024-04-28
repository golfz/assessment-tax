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

func TestSetPersonalDeductionHandler_Integration_WithDefault_Success(t *testing.T) {
	// Arrange
	input := admin.Deduction{
		Deduction: deduction.DefaultPersonalDeduction,
	}
	want := admin.PersonalDeduction{
		Deduction: deduction.DefaultPersonalDeduction,
	}

	cfg := config.NewWith(os.Getenv)
	pg, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		t.Errorf("failed to connect to database: %v", err)
	}

	e := echo.New()
	hAdmin := admin.New(pg)
	e.POST("/admin/deductions/personal", hAdmin.SetPersonalDeductionHandler)

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

	var got admin.PersonalDeduction
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
	}
	assert.Equal(t, want, got)

}

func TestSetKReceiptDeductionHandler_Integration_WithDefault_Success(t *testing.T) {
	// Arrange
	input := admin.Deduction{
		Deduction: deduction.DefaultKReceiptDeduction,
	}
	want := admin.KReceiptDeduction{
		Deduction: deduction.DefaultKReceiptDeduction,
	}

	cfg := config.NewWith(os.Getenv)
	pg, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		t.Errorf("failed to connect to database: %v", err)
	}

	e := echo.New()
	hAdmin := admin.New(pg)
	e.POST("/admin/deductions/k-receipt", hAdmin.SetKReceiptDeductionHandler)

	var bReader io.Reader
	b, _ := json.Marshal(input)
	bReader = strings.NewReader(string(b))
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var got admin.KReceiptDeduction
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
	}
	assert.Equal(t, want, got)
}
