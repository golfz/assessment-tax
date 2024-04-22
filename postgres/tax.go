package postgres

import (
	"github.com/golfz/assessment-tax/tax"
	"log"
)

func (p *Postgres) GetDeduction() (tax.Deduction, error) {
	selectSql := `SELECT name, amount FROM deductions`
	rows, err := p.Db.Query(selectSql)
	if err != nil {
		log.Printf("unable to query deduction: %v", err)
		return tax.Deduction{}, err
	}
	defer rows.Close()

	var deduction tax.Deduction
	for rows.Next() {
		var name string
		var amount float64
		err = rows.Scan(&name, &amount)
		if err != nil {
			log.Printf("unable to scan deduction: %v", err)
			return tax.Deduction{}, err
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
