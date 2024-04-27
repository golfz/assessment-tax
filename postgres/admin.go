package postgres

func (p *Postgres) SetPersonalDeduction(amount float64) error {
	updateSql := `UPDATE deductions SET amount = $1 WHERE name = 'personal'`
	_, err := p.DB.Exec(updateSql, amount)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) SetKReceiptDeduction(amount float64) error {
	updateSql := `UPDATE deductions SET amount = $1 WHERE name = 'k-receipt'`
	_, err := p.DB.Exec(updateSql, amount)
	if err != nil {
		return err
	}

	return nil
}
