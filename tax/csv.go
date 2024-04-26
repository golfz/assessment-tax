package tax

import (
	"encoding/csv"
	"io"
	"strconv"
)

const (
	csvHeaderRowIndex = 0
)

type csvColumn int

const (
	csvColumnTotalIncome csvColumn = iota
	csvColumnWHT
	csvColumnDonation
)

func validateHeader(header []string) error {
	if len(header) != 3 {
		return ErrInvalidCSVHeader
	}

	if header[0] != "totalIncome" || header[1] != "wht" || header[2] != "donation" {
		return ErrInvalidCSVHeader
	}

	return nil
}

func getColumnValue(row []string, column csvColumn) (float64, error) {
	result, err := strconv.ParseFloat(row[column], 64)
	if err != nil {
		return 0, ErrParsingData
	}
	return result, nil
}

func getTaxInformation(row []string) (TaxInformation, error) {
	taxInfo := TaxInformation{}

	if len(row) != 3 {
		return TaxInformation{}, ErrParsingData
	}

	var err error

	if taxInfo.TotalIncome, err = getColumnValue(row, csvColumnTotalIncome); err != nil {
		return TaxInformation{}, err
	}

	if taxInfo.WHT, err = getColumnValue(row, csvColumnWHT); err != nil {
		return TaxInformation{}, err
	}

	donation, err := getColumnValue(row, csvColumnDonation)
	if err != nil {
		return TaxInformation{}, err
	}
	taxInfo.Allowances = []Allowance{
		{
			Type:   AllowanceTypeDonation,
			Amount: donation,
		},
	}

	return taxInfo, nil
}

func parseRow(index int, row []string, data *[]TaxInformation) error {
	if index == csvHeaderRowIndex {
		return validateHeader(row)
	}

	taxInfo, err := getTaxInformation(row)
	if err != nil {
		return err
	}

	*data = append(*data, taxInfo)
	return nil
}

func parseTaxRecords(records [][]string) ([]TaxInformation, error) {
	result := make([]TaxInformation, 0)

	for rowIndex, rowData := range records {
		err := parseRow(rowIndex, rowData, &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func readCSV(r io.Reader) ([]TaxInformation, error) {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, ErrReadingCSV
	}

	return parseTaxRecords(records)
}
