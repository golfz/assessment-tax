package postgres

import (
	"errors"
	"github.com/golfz/assessment-tax/deduction"
)

var (
	ErrCannotQueryDeduction = errors.New("unable to query deduction")
	ErrCannotScanDeduction  = errors.New("unable to scan deduction")
)

func (p *Postgres) GetDeduction() (deduction.Deduction, error) {
	selectSql := `SELECT name, amount FROM deductions`
	rows, err := p.Db.Query(selectSql)
	if err != nil {
		return deduction.Deduction{}, ErrCannotQueryDeduction
	}
	defer rows.Close()

	var deductionData deduction.Deduction
	for rows.Next() {
		var name string
		var amount float64
		err = rows.Scan(&name, &amount)
		if err != nil {
			return deduction.Deduction{}, ErrCannotScanDeduction
		}
		switch name {
		case "personal":
			deductionData.Personal = amount
		case "k-receipt":
			deductionData.KReceipt = amount
		case "donation":
			deductionData.Donation = amount
		}
	}

	return deductionData, nil
}
