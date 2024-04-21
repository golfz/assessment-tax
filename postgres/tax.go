package postgres

import "github.com/golfz/assessment-tax/tax"

func (p *Postgres) GetDeduction() (tax.Deduction, error) {
	return tax.Deduction{Personal: 60_000.0}, nil
}
