package admin

import (
	"encoding/json"
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
	// Arrange
	amount := 70_000.0
	rec, c, h, mock := setup(http.MethodPost, "/admin/deductions/personal", Input{Amount: amount})
	mock.ExpectToCall(MethodSetPersonalDeduction)

	// Act
	err := h.SetPersonalDeductionHandler(c)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
	mock.Verify(t)
	assert.Equal(t, amount, mock.whatIsAmount)
}
