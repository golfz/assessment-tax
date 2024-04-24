package postgres

import (
	"errors"
	"github.com/golfz/assessment-tax/tax"
)

var (
	ErrCannotQueryDeduction = errors.New("unable to query deduction")
	ErrCannotScanDeduction  = errors.New("unable to scan deduction")
)

func (p *Postgres) GetDeduction() (tax.Deduction, error) {
	selectSql := `SELECT name, amount FROM deductions`
	rows, err := p.Db.Query(selectSql)
	if err != nil {
		return tax.Deduction{}, ErrCannotQueryDeduction
	}
	defer rows.Close()

	var deduction tax.Deduction
	for rows.Next() {
		var name string
		var amount float64
		err = rows.Scan(&name, &amount)
		if err != nil {
			return tax.Deduction{}, ErrCannotScanDeduction
		}
		switch name {
		case "personal":
			deduction.Personal = amount
		case "k-receipt":
			deduction.KReceipt = amount
		case "donation":
			deduction.Donation = amount
		}
	}

	return deduction, nil
}
