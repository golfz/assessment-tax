package admin

type Input struct {
	Amount float64 `json:"amount" validate:"min=0"`
}

type PersonalDeduction struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type KReceiptDeduction struct {
	KReceiptDeduction float64 `json:"kReceipt"`
}
