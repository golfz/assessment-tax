package postgres

func (p *Postgres) SetPersonalDeduction(amount float64) error {
	updateSql := `UPDATE deductions SET amount = $1 WHERE name = 'personal'`
	_, err := p.Db.Exec(updateSql, amount)
	if err != nil {
		return err
	}

	return nil
}
