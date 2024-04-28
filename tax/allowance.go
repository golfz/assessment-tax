package tax

import (
	"github.com/golfz/assessment-tax/deduction"
	"math"
)

func collapseAllowance(allowances []Allowance) map[AllowanceType]float64 {
	result := make(map[AllowanceType]float64)
	for _, a := range allowances {
		result[a.Type] += a.Amount
	}
	return result
}

func getTaxableAllowance(allowances []Allowance, deduction deduction.Deduction) map[AllowanceType]float64 {
	result := collapseAllowance(allowances)

	for aType, aAmount := range result {
		switch aType {
		case AllowanceTypeDonation:
			result[aType] = math.Min(aAmount, deduction.Donation)
		case AllowanceTypeKReceipt:
			result[aType] = math.Min(aAmount, deduction.KReceipt)
		}
	}
	return result
}

func getTotalAllowance(allowances []Allowance, deduction deduction.Deduction) float64 {
	taxableAllowances := getTaxableAllowance(allowances, deduction)
	total := 0.0
	for _, aAmount := range taxableAllowances {
		total += aAmount
	}
	return total
}
