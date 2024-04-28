package admin

type Deduction struct {
	Deduction float64 `json:"amount" validate:"min=0"`
}

type PersonalDeduction struct {
	Deduction float64 `json:"personalDeduction"`
}

type KReceiptDeduction struct {
	Deduction float64 `json:"kReceipt"`
}
