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

type Deduction struct {
	Personal float64
}

func CalculateTax(info TaxInformation, deduction Deduction) (TaxResult, error) {
	netIncome := info.TotalIncome - deduction.Personal

	// Calculate tax
	tax := 0.0

	// 0 - 150,000 = 0%
	left := 0.0
	right := 150_000.0
	taxRate := 0.0
	if netIncome > left {
		taxableIncome := netIncome - left
		taxRange := right - left
		if taxableIncome > taxRange {
			taxableIncome = right - left
		}
		tax += taxableIncome * (taxRate / 100.0)
	}

	// 150,001 - 500,000 = 10%
	left = 150_000.0
	right = 500_000.0
	taxRate = 10.00
	if netIncome > left {
		taxableIncome := netIncome - left
		taxRange := right - left
		if taxableIncome > taxRange {
			taxableIncome = right - left
		}
		tax += taxableIncome * (taxRate / 100.0)
	}

	// 500,001 - 1,000,000 = 15%
	left = 500_000.0
	right = 1_000_000.0
	taxRate = 15.00
	if netIncome > left {
		taxableIncome := netIncome - left
		taxRange := right - left
		if taxableIncome > taxRange {
			taxableIncome = right - left
		}
		tax += taxableIncome * (taxRate / 100.0)
	}

	return TaxResult{
		Tax: tax,
	}, nil
}
