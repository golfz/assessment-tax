package tax

import "math"

func collapseAllowance(allowances []Allowance) map[AllowanceType]float64 {
	result := make(map[AllowanceType]float64)
	for _, a := range allowances {
		result[a.Type] += a.Amount // Directly add 'a.Amount' to the map value.
	}
	return result
}

func getTaxableAllowance(allowances []Allowance, deduction Deduction) map[AllowanceType]float64 {
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
