package tax

import (
	"encoding/csv"
	"io"
	"strconv"
)

type CSVReader struct {
	reader io.Reader
}

func NewCSVReader(r io.Reader) *CSVReader {
	return &CSVReader{reader: r}
}

const (
	csvHeaderRowIndex = 0
)

type csvColumn int

const (
	csvColumnTotalIncome csvColumn = iota
	csvColumnWHT
	csvColumnDonation
)

func (cr *CSVReader) validateHeader(header []string) error {
	if len(header) != 3 {
		return ErrInvalidCSVHeader
	}

	if header[0] != "totalIncome" || header[1] != "wht" || header[2] != "donation" {
		return ErrInvalidCSVHeader
	}

	return nil
}

func (cr *CSVReader) getColumnValue(row []string, column csvColumn) (float64, error) {
	result, err := strconv.ParseFloat(row[column], 64)
	if err != nil {
		return 0, ErrParsingData
	}
	return result, nil
}

func (cr *CSVReader) getTaxInformation(row []string) (TaxInformation, error) {
	taxInfo := TaxInformation{}

	if len(row) != 3 {
		return TaxInformation{}, ErrParsingData
	}

	var err error

	if taxInfo.TotalIncome, err = cr.getColumnValue(row, csvColumnTotalIncome); err != nil {
		return TaxInformation{}, err
	}

	if taxInfo.WHT, err = cr.getColumnValue(row, csvColumnWHT); err != nil {
		return TaxInformation{}, err
	}

	donation, err := cr.getColumnValue(row, csvColumnDonation)
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

func (cr *CSVReader) parseRow(index int, row []string, data *[]TaxInformation) error {
	if index == csvHeaderRowIndex {
		return cr.validateHeader(row)
	}

	taxInfo, err := cr.getTaxInformation(row)
	if err != nil {
		return err
	}

	*data = append(*data, taxInfo)
	return nil
}

func (cr *CSVReader) parseTaxRecords(records [][]string) ([]TaxInformation, error) {
	result := make([]TaxInformation, 0)

	for rowIndex, rowData := range records {
		err := cr.parseRow(rowIndex, rowData, &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (cr *CSVReader) readRecords() ([]TaxInformation, error) {
	reader := csv.NewReader(cr.reader)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, ErrReadingCSV
	}

	return cr.parseTaxRecords(records)
}
