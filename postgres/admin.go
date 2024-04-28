package postgres

func (p *Postgres) setDeduction(deductionType string, amount float64) error {
	updateSql := "UPDATE deductions SET amount = $1 WHERE name = $2"
	_, err := p.DB.Exec(updateSql, amount, deductionType)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) SetPersonalDeduction(amount float64) error {
	return p.setDeduction("personal", amount)
}

func (p *Postgres) SetKReceiptDeduction(amount float64) error {
	return p.setDeduction("k-receipt", amount)
}
