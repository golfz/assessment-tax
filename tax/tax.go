package tax

type AllowanceType string

const (
	AllowanceTypeDonation AllowanceType = "donation"
	AllowanceTypeKReceipt AllowanceType = "k-receipt"
)

type Allowance struct {
	Type   AllowanceType
	Amount float64
}

type TaxInformation struct {
	TotalIncome float64
	WHT         float64
	Allowances  []Allowance
}

type TaxResult struct {
	Tax float64
}

func CalculateTax(info TaxInformation) (TaxResult, error) {
	personalDeduction := 60_000.0
	taxableIncome := info.TotalIncome - personalDeduction

	// Calculate tax
	tax := 0.0

	// 0 - 150,000 = 0%
	left := 0.0
	right := 150_000.0
	taxRate := 0.0
	if taxableIncome > left {
		incomeForThisRate := taxableIncome - left
		if incomeForThisRate > right {
			incomeForThisRate = right - left
		}
		tax += incomeForThisRate * (taxRate / 100.0)
	}

	// 150,001 - 500,000 = 10%
	left = 150_000.0
	right = 500_000.0
	taxRate = 10.00
	if taxableIncome > left {
		incomeForThisRate := taxableIncome - left
		if incomeForThisRate > right {
			incomeForThisRate = right - left
		}
		tax += incomeForThisRate * (taxRate / 100.0)
	}

	return TaxResult{
		Tax: tax,
	}, nil
}
