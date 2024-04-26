package tax

type AllowanceType string

const (
	AllowanceTypeDonation AllowanceType = "donation"
	AllowanceTypeKReceipt AllowanceType = "k-receipt"
)

type Allowance struct {
	Type   AllowanceType `json:"allowanceType"`
	Amount float64       `json:"amount" validate:"min=0"`
}

type TaxInformation struct {
	TotalIncome float64     `json:"totalIncome" validate:"required,min=0"`
	WHT         float64     `json:"wht" validate:"min=0"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResult struct {
	Tax       float64    `json:"tax"`
	TaxRefund float64    `json:"taxRefund,omitempty"`
	TaxLevels []TaxLevel `json:"taxLevel"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type rate struct {
	lowerBound  float64
	upperBound  float64
	percentage  float64
	description string
}

type CsvTaxRequest struct {
	TotalIncome float64 `csv:"totalIncome"`
	WHT         float64 `csv:"wht"`
	Donation    float64 `csv:"donation"`
}

type CsvTaxResponse struct {
	Taxes []CsvTaxRecord `json:"taxes"`
}

type CsvTaxRecord struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
	TaxRefund   float64 `json:"taxRefund,omitempty"`
}
