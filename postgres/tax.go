package postgres

import (
	"errors"
	"github.com/golfz/assessment-tax/deduction"
)

var (
	ErrCannotQueryDeduction = errors.New("unable to query deduction")
	ErrCannotScanDeduction  = errors.New("unable to scan deduction")
)

const (
	namePersonalDeduction = "personal"
	nameKReceiptDeduction = "k-receipt"
	nameDonationDeduction = "donation"
)

func applyDeductionValue(name string, amount float64, deductionData *deduction.Deduction) {
	switch name {
	case namePersonalDeduction:
		deductionData.Personal = amount
	case nameKReceiptDeduction:
		deductionData.KReceipt = amount
	case nameDonationDeduction:
		deductionData.Donation = amount
	}
}

func (p *Postgres) GetDeduction() (deduction.Deduction, error) {
	selectSQL := `SELECT name, amount FROM deductions`
	rows, err := p.DB.Query(selectSQL)
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

		applyDeductionValue(name, amount, &deductionData)
	}

	return deductionData, nil
}
