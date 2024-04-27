//go:build unit

package admin

import (
	"encoding/json"
	"github.com/golfz/assessment-tax/deduction"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	MethodSetPersonalDeduction = "SetPersonalDeduction"
	MethodSetKReceiptDeduction = "SetKReceiptDeduction"
)

type mockAdminStorer struct {
	err          error
	methodToCall map[string]bool
	whatIsAmount float64
}

func NewMockTaxStorer() *mockAdminStorer {
	return &mockAdminStorer{
		methodToCall: make(map[string]bool),
	}
}

func (m *mockAdminStorer) SetPersonalDeduction(amount float64) error {
	m.methodToCall[MethodSetPersonalDeduction] = true
	m.whatIsAmount = amount
	return m.err
}

func (m *mockAdminStorer) SetKReceiptDeduction(amount float64) error {
	m.methodToCall[MethodSetKReceiptDeduction] = true
	m.whatIsAmount = amount
	return m.err
}

func (m *mockAdminStorer) ExpectToCall(methodName string) {
	if m.methodToCall == nil {
		m.methodToCall = make(map[string]bool)
	}
	m.methodToCall[methodName] = false
}

func (m *mockAdminStorer) Verify(t *testing.T) {
	for methodName, called := range m.methodToCall {
		if !called {
			t.Errorf("expected %s to be called", methodName)
		}
	}
}

func setup(method, url string, body interface{}) (*httptest.ResponseRecorder, echo.Context, *Handler, *mockAdminStorer) {
	var bReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bReader = strings.NewReader(string(b))
	}
	req := httptest.NewRequest(method, url, bReader)
	if body != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	mock := NewMockTaxStorer()
	h := New(mock)

	return rec, c, h, mock
}

func TestSetPersonalDeductionHandler_Success(t *testing.T) {
	testCases := []struct {
		name   string
		amount float64
	}{
		{
			name:   "EXP05: setting personal deduction",
			amount: 70_000.0,
		},
		{
			name:   "setting with minimum personal deduction",
			amount: deduction.MinPersonalDeduction + 1,
		},
		{
			name:   "setting with maximum personal deduction",
			amount: deduction.MaxPersonalDeduction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			rec, c, h, mock := setup(http.MethodPost, "/admin/deductions/personal", Input{Amount: tc.amount})
			mock.ExpectToCall(MethodSetPersonalDeduction)

			// Act
			err := h.SetPersonalDeductionHandler(c)

			// Assert
			mock.Verify(t)
			assert.Equal(t, tc.amount, mock.whatIsAmount)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			var got PersonalDeduction
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
			}
			assert.Equal(t, tc.amount, got.PersonalDeduction)
		})
	}
}

func TestSetPersonalDeductionHandler_Error(t *testing.T) {
	t.Run("no content-type header", func(t *testing.T) {
		// Arrange
		body := struct{ Field string }{Field: "invalid"}
		resp, c, h, _ := setup(http.MethodPost, "/admin/deductions/personal", body)
		c.Request().Header.Set(echo.HeaderContentType, "")

		// Act
		err := h.SetPersonalDeductionHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
		assert.Equal(t, ErrReadingRequestBody.Error(), got.Message)
	})

	t.Run("invalid input", func(t *testing.T) {
		// Arrange
		body := struct{ Amount float64 }{Amount: -1}
		resp, c, h, _ := setup(http.MethodPost, "/admin/deductions/personal", body)

		// Act
		err := h.SetPersonalDeductionHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
		assert.Equal(t, ErrInvalidInput.Error(), got.Message)
	})

	t.Run("SetPersonalDeduction() error", func(t *testing.T) {
		// Arrange
		body := struct{ Amount float64 }{Amount: 70_000}
		resp, c, h, mock := setup(http.MethodPost, "/admin/deductions/personal", body)
		mock.err = ErrSettingPersonalDeduction

		// Act
		err := h.SetPersonalDeductionHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
		assert.Equal(t, ErrSettingPersonalDeduction.Error(), got.Message)
	})
}

func TestSetPersonalDeductionHandler_ValidateAmount_Error(t *testing.T) {
	testCases := []struct {
		name   string
		amount float64
	}{
		{
			name:   "amount less than minimum personal deduction; expected error",
			amount: deduction.MinPersonalDeduction - 1,
		},
		{
			name:   "amount equal minimum personal deduction boundary; expected error",
			amount: deduction.MinPersonalDeduction,
		},
		{
			name:   "amount more than maximum personal deduction",
			amount: deduction.MaxPersonalDeduction + 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			rec, c, h, _ := setup(http.MethodPost, "/admin/deductions/personal", Input{Amount: tc.amount})

			// Act
			err := h.SetPersonalDeductionHandler(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			var got Err
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
			}
			assert.NotEmpty(t, got.Message)
			assert.Equal(t, ErrInvalidPersonalDeduction.Error(), got.Message)
		})
	}
}

func TestSetKReceiptDeductionHandler_Success(t *testing.T) {
	testCases := []struct {
		name   string
		amount float64
	}{
		{
			name:   "EXP08: setting k-receipt deduction",
			amount: 70_000.0,
		},
		{
			name:   "setting with minimum k-receipt deduction",
			amount: deduction.MinKReceiptDeduction + 1,
		},
		{
			name:   "setting with maximum k-receipt deduction",
			amount: deduction.MaxKReceiptDeduction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			rec, c, h, mock := setup(http.MethodPost, "/admin/deductions/k-receipt", Input{Amount: tc.amount})
			mock.ExpectToCall(MethodSetKReceiptDeduction)

			// Act
			err := h.SetKReceiptDeductionHandler(c)

			// Assert
			mock.Verify(t)
			assert.Equal(t, tc.amount, mock.whatIsAmount)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			var got KReceiptDeduction
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
			}
			assert.Equal(t, tc.amount, got.KReceiptDeduction)
		})
	}
}

func TestSetKReceiptDeductionHandler_Error(t *testing.T) {
	t.Run("no content-type header", func(t *testing.T) {
		// Arrange
		body := struct{ Field string }{Field: "invalid"}
		resp, c, h, _ := setup(http.MethodPost, "/admin/deductions/k-receipt", body)
		c.Request().Header.Set(echo.HeaderContentType, "")

		// Act
		err := h.SetKReceiptDeductionHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
		assert.Equal(t, ErrReadingRequestBody.Error(), got.Message)
	})
}
