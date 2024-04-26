//go:build unit

package tax

import (
	"github.com/golfz/assessment-tax/deduction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollapseAllowance_WithEmptyAllowance_ExpectEmpty(t *testing.T) {
	// Arrange
	allowances := []Allowance{}

	// Act
	result := collapseAllowance(allowances)

	// Assert
	assert.Empty(t, result)
}

func TestCollapseAllowance_WithSingleAllowance_ExpectSameAllowance(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 100_000.0},
	}

	// Act
	result := collapseAllowance(allowances)

	// Assert
	assert.Equal(t, map[AllowanceType]float64{
		AllowanceTypeDonation: 100_000.0,
	}, result)
}

func TestCollapseAllowance_WithMultipleAllowance_ExpectSummedAllowance(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 100_000.0},
		{Type: AllowanceTypeDonation, Amount: 50_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 50_000.0},
	}

	// Act
	result := collapseAllowance(allowances)

	// Assert
	assert.Equal(t, map[AllowanceType]float64{
		AllowanceTypeDonation: 150_000.0,
		AllowanceTypeKReceipt: 50_000.0,
	}, result)
}

func TestGetTaxableAllowance_WithEmptyAllowance_ExpectEmpty(t *testing.T) {
	// Arrange
	allowances := []Allowance{}
	defaultDeduction := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTaxableAllowance(allowances, defaultDeduction)

	// Assert
	assert.Empty(t, result)
}

func TestGetTaxableAllowance_WithSingleAllowance_ExpectSameAllowance(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 80_000.0},
	}
	defaultDeduction := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTaxableAllowance(allowances, defaultDeduction)

	// Assert
	assert.Equal(t, map[AllowanceType]float64{
		AllowanceTypeDonation: 80_000.0,
	}, result)
}

func TestGetTaxableAllowance_WithMultipleAllowance_ExpectTaxableAllowance(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 30_000.0},
		{Type: AllowanceTypeDonation, Amount: 40_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 50_000.0},
	}
	defaultDeduction := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTaxableAllowance(allowances, defaultDeduction)

	// Assert
	assert.Equal(t, map[AllowanceType]float64{
		AllowanceTypeDonation: 70_000.0,
		AllowanceTypeKReceipt: 50_000.0,
	}, result)
}

func TestGetTaxableAllowance_WithAllowanceMoreThanDeduction_ExpectSameDeduction(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 130_000.0},
		{Type: AllowanceTypeDonation, Amount: 70_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 200_000.0},
	}
	defaultDeduction := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTaxableAllowance(allowances, defaultDeduction)

	// Assert
	assert.Equal(t, map[AllowanceType]float64{
		AllowanceTypeDonation: defaultDeduction.Donation,
		AllowanceTypeKReceipt: defaultDeduction.KReceipt,
	}, result)
}

func TestGetTaxableAllowance_WithAllowanceEqualDeduction_ExpectSameAllowance(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 40_000.0},
		{Type: AllowanceTypeDonation, Amount: 60_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 20_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 30_000.0},
	}
	defaultDeduction := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTaxableAllowance(allowances, defaultDeduction)

	// Assert
	assert.Equal(t, map[AllowanceType]float64{
		AllowanceTypeDonation: 100_000.0,
		AllowanceTypeKReceipt: 50_000.0,
	}, result)
}

func TestGetTotalAllowance_WithEmptyAllowance_ExpectZero(t *testing.T) {
	// Arrange
	allowances := []Allowance{}
	deductionData := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTotalAllowance(allowances, deductionData)

	// Assert
	assert.Equal(t, 0.0, result)
}

func TestGetTotalAllowance_WithSingleAllowance_ExpectSameAllowance(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 80_000.0},
	}
	deductionData := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTotalAllowance(allowances, deductionData)

	// Assert
	assert.Equal(t, 80_000.0, result)
}

func TestGetTotalAllowance_WithMultipleAllowance_ExpectTotalAllowance(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 30_000.0},
		{Type: AllowanceTypeDonation, Amount: 40_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 50_000.0},
	}
	deductionData := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTotalAllowance(allowances, deductionData)

	// Assert
	assert.Equal(t, 120_000.0, result)
}

func TestGetTotalAllowance_WithAllowanceMoreThanDeduction_ExpectTotalDeduction(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 130_000.0},
		{Type: AllowanceTypeDonation, Amount: 70_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 200_000.0},
	}
	deductionData := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTotalAllowance(allowances, deductionData)

	// Assert
	assert.Equal(t, 150_000.0, result)
}

func TestGetTotalAllowance_WithAllowanceEqualDeduction_ExpectTotalAllowance(t *testing.T) {
	// Arrange
	allowances := []Allowance{
		{Type: AllowanceTypeDonation, Amount: 40_000.0},
		{Type: AllowanceTypeDonation, Amount: 60_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 20_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 30_000.0},
	}
	deductionData := deduction.Deduction{
		Personal: 60_000.0,
		KReceipt: 50_000.0,
		Donation: 100_000.0,
	}

	// Act
	result := getTotalAllowance(allowances, deductionData)

	// Assert
	assert.Equal(t, 150_000.0, result)
}
