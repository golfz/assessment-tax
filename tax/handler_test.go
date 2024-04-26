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
		name          string
		taxInfo       TaxInformation
		wantTaxResult TaxResult
	}{
		{
			name: "EXP01: basic income, no WHT, no Allowance; expect tax",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 29_000.0, TaxRefund: 0.0},
		},
		{
			name: "EXP02: Income and WHT, no Allowance; expect tax",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         25_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 0.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 4_000.0, TaxRefund: 0.0},
		},
		{
			name: "EXP03: Income and Allowance, no WHT; expect tax",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 19_000.0, TaxRefund: 0.0},
		},
		{
			name: "One Allowance, tax payable > WHT; expect tax",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         15_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 4_000.0, TaxRefund: 0.0},
		},
		{
			name: "One Allowance, tax payable = WHT; expect tax=0",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         19_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 0.0, TaxRefund: 0.0},
		},
		{
			name: "One Allowance, tax payable < WHT; expect taxRefund",
			taxInfo: TaxInformation{
				TotalIncome: 500_000.0,
				WHT:         29_000.0,
				Allowances: []Allowance{
					{Type: AllowanceTypeDonation, Amount: 200_000.0},
				},
			},
			wantTaxResult: TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
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
			wantTaxResult: TaxResult{Tax: 9_000.0, TaxRefund: 0.0},
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
			wantTaxResult: TaxResult{Tax: 0.0, TaxRefund: 0.0},
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
			wantTaxResult: TaxResult{Tax: 0.0, TaxRefund: 10_000.0},
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
			assert.Equal(t, tc.wantTaxResult.Tax, got.Tax)
			assert.Equal(t, tc.wantTaxResult.TaxRefund, got.TaxRefund)
		})
	}
}

func TestCalculateTaxHandler_WithTaxLevel_Success(t *testing.T) {
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
			assert.Equal(t, tc.wantTaxResult.Tax, got.Tax)
			assert.Equal(t, tc.wantTaxResult.TaxRefund, got.TaxRefund)
			for i, wantTax := range tc.wantTaxLevels {
				assert.Equal(t, wantTax, got.TaxLevels[i].Tax)
			}
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
		taxInfo := TaxInformation{
			TotalIncome: 500_000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{Type: AllowanceTypeDonation, Amount: 0.0},
			},
		}
		resp, c, h, mock := setup(http.MethodPost, "/tax/calculations", taxInfo)
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
		taxInfo := TaxInformation{
			TotalIncome: 500_000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{Type: AllowanceTypeDonation, Amount: 0.0},
			},
		}
		resp, c, h, mock := setup(http.MethodPost, "/tax/calculations", taxInfo)
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
		taxInfo := TaxInformation{
			TotalIncome: -1,
			WHT:         0.0,
			Allowances: []Allowance{
				{Type: AllowanceTypeDonation, Amount: 0.0},
			},
		}
		resp, c, h, _ := setup(http.MethodPost, "/tax/calculations", taxInfo)

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
		taxInfo := TaxInformation{
			TotalIncome: 100_000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{Type: AllowanceTypeDonation, Amount: -10.0},
			},
		}
		resp, c, h, _ := setup(http.MethodPost, "/tax/calculations", taxInfo)

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
