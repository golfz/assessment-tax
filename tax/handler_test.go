//go:build unit

package tax

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	MethodGetDeduction = "GetDeduction"
)

type mockTaxStorer struct {
	result       TaxResult
	deduction    Deduction
	err          error
	methodToCall map[string]bool
}

func NewMockTaxStorer() *mockTaxStorer {
	return &mockTaxStorer{
		methodToCall: make(map[string]bool),
	}
}

func (m *mockTaxStorer) GetDeduction() (Deduction, error) {
	m.methodToCall[MethodGetDeduction] = true
	return m.deduction, m.err
}

func (m *mockTaxStorer) ExpectToCall(methodName string) {
	if m.methodToCall == nil {
		m.methodToCall = make(map[string]bool)
	}
	m.methodToCall[methodName] = false
}

func (m *mockTaxStorer) Verify(t *testing.T) {
	for methodName, called := range m.methodToCall {
		if !called {
			t.Errorf("expected %s to be called", methodName)
		}
	}
}

func setup(method, url string, body interface{}) (*httptest.ResponseRecorder, echo.Context, *Handler, *mockTaxStorer) {
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

func TestCalculateTaxHandler_Success(t *testing.T) {
	deduction := Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	testcases := []struct {
		name       string
		taxInfo    TaxInformation
		wantTax    float64
		wantRefund float64
	}{
		{
			name: "EXP01: expect 200 OK with tax",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTax:    29_000.0,
			wantRefund: 0.0,
		},
		{
			name: "EXP02: tax > WHT; expect 200 OK with tax>0",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         25_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTax:    4_000.0,
			wantRefund: 0.0,
		},
		{
			name: "tax = WHT; expect 200 OK with tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         29_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTax:    0.0,
			wantRefund: 0.0,
		},
		{
			name: "tax < WHT; expect 200 OK with tax=0, taxRefund>0",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         39_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTax:    0.0,
			wantRefund: 10_000.0,
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
			wantTax:    19_000.0,
			wantRefund: 0.0,
		},
		{
			name: "income=500,000 wht=tax donation=200,000; expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         19_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTax:    0.0,
			wantRefund: 0.0,
		},
		{
			name: "income=500,000 wht>tax donation=200,000; expect taxRefund=10,000",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         29_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTax:    0.0,
			wantRefund: 10_000.0,
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
			wantTax:    0.0,
			wantRefund: 10_000.0,
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
			wantTax:    0.0,
			wantRefund: 10_000.0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			resp, c, h, mock := setup(http.MethodPost, "/tax/calculations", tc.taxInfo)
			mock.err = nil
			mock.deduction = deduction
			mock.ExpectToCall(MethodGetDeduction)

			// Act
			err := h.CalculateTaxHandler(c)

			// Assert
			mock.Verify(t)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
			var got TaxResult
			if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
				t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
			}
			assert.Equal(t, tc.wantTax, got.Tax)
			assert.Equal(t, tc.wantRefund, got.TaxRefund)
		})
	}
}

func TestCalculateTaxHandler_Error(t *testing.T) {
	t.Run("no content-type expect 400 with error message", func(t *testing.T) {
		// Arrange
		body := struct{ Field string }{Field: "invalid"}
		resp, c, h, _ := setup(http.MethodPost, "/tax/calculations", body)
		c.Request().Header.Set(echo.HeaderContentType, "")

		// Act
		err := h.CalculateTaxHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
		assert.Equal(t, "cannot reading request body", got.Message)
	})

	t.Run("incorrect body expect 400 with error message", func(t *testing.T) {
		// Arrange
		body := struct{ Field string }{Field: "invalid"}
		resp, c, h, _ := setup(http.MethodPost, "/tax/calculations", body)

		// Act
		err := h.CalculateTaxHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
		assert.Equal(t, ErrInvalidTaxInformation.Error(), got.Message)
	})

	t.Run("GetDeduction() error expect 500 with error message", func(t *testing.T) {
		// Arrange
		info := TaxInformation{
			TotalIncome: 500_000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{Type: AllowanceTypeDonation, Amount: 0.0},
			},
		}
		resp, c, h, mock := setup(http.MethodPost, "/tax/calculations", info)
		mock.err = errors.New("error getting deduction")
		mock.ExpectToCall(MethodGetDeduction)

		// Act
		err := h.CalculateTaxHandler(c)

		// Assert
		mock.Verify(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.Equal(t, "error getting deduction", got.Message)
	})

	t.Run("invalid deduction expect 500 with error message", func(t *testing.T) {
		// Arrange
		info := TaxInformation{
			TotalIncome: 500_000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{Type: AllowanceTypeDonation, Amount: 0.0},
			},
		}
		resp, c, h, mock := setup(http.MethodPost, "/tax/calculations", info)
		mock.deduction = Deduction{}
		mock.ExpectToCall(MethodGetDeduction)

		// Act
		err := h.CalculateTaxHandler(c)

		// Assert
		mock.Verify(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.Equal(t, "error calculating tax", got.Message)
	})

	t.Run("invalid total income expect 400 with error message", func(t *testing.T) {
		// Arrange
		info := TaxInformation{
			TotalIncome: -1,
			WHT:         0.0,
			Allowances: []Allowance{
				{Type: AllowanceTypeDonation, Amount: 0.0},
			},
		}
		resp, c, h, _ := setup(http.MethodPost, "/tax/calculations", info)

		// Act
		err := h.CalculateTaxHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.Equal(t, ErrInvalidTaxInformation.Error(), got.Message)
	})

	t.Run("invalid allowance amount expect 400 with error message", func(t *testing.T) {
		// Arrange
		info := TaxInformation{
			TotalIncome: 100_000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{Type: AllowanceTypeDonation, Amount: -10.0},
			},
		}
		resp, c, h, _ := setup(http.MethodPost, "/tax/calculations", info)

		// Act
		err := h.CalculateTaxHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.Equal(t, ErrInvalidTaxInformation.Error(), got.Message)
	})
}
