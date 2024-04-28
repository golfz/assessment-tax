package postgres

type deductionType string

const (
	personalDeduction  deductionType = "personal"
	kReceiptDeduction  deductionType = "k-receipt"
	updateDeductionSQL               = "UPDATE deductions SET amount = $1 WHERE name = $2"
)

func (p *Postgres) setDeduction(deducType deductionType, amount float64) error {
	_, err := p.DB.Exec(updateDeductionSQL, amount, deducType)
	return err
}

func (p *Postgres) SetPersonalDeduction(amount float64) error {
	return p.setDeduction(personalDeduction, amount)
}

func (p *Postgres) SetKReceiptDeduction(amount float64) error {
	return p.setDeduction(kReceiptDeduction, amount)
}
