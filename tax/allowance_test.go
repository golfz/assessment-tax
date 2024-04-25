//go:build unit

package tax

import (
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
	assert.Equal(t, allowances, result)
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
	assert.Equal(t, []Allowance{
		{Type: AllowanceTypeDonation, Amount: 150_000.0},
		{Type: AllowanceTypeKReceipt, Amount: 50_000.0},
	}, result)
}
