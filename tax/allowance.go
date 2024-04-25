package tax

func collapseAllowance(allowances []Allowance) []Allowance {
	allowanceSums := make(map[AllowanceType]float64)
	for _, a := range allowances {
		allowanceSums[a.Type] += a.Amount // Directly add 'a.Amount' to the map value.
	}
	result := make([]Allowance, 0, len(allowanceSums))
	for k, v := range allowanceSums {
		result = append(result, Allowance{Type: k, Amount: v})
	}
	return result
}
