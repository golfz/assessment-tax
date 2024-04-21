package tax

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCalculateTax_onlyTotalIncome(t *testing.T) {
	// Arrange
	info := TaxInformation{
		TotalIncome: 500_000.0,
		WHT:         0.0,
		Allowances: []Allowance{
			{
				Type:   AllowanceTypeDonation,
				Amount: 0.0,
			},
		},
	}
	deduction := Deduction{Personal: 60_000.0}
	want := TaxResult{
		Tax: 29_000.0,
	}

	// Act
	got, err := CalculateTax(info, deduction)

	// Assert
	fmt.Printf("got: %#v\n", got)
	if err != nil {
		t.Errorf("CalculateTax() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CalculateTax() = %v, want %v", got, want)
	}
}
