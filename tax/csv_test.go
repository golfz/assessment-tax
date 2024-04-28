//go:build unit

package tax

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTaxRecords_Success(t *testing.T) {
	// Arrange
	testCases := []struct {
		name    string
		records [][]string
		want    []TaxInformation
	}{
		{
			name: "multiple records",
			records: [][]string{
				{"totalIncome", "wht", "donation"},
				{"1000000", "100000", "10000"},
				{"2000000", "200000", "20000"},
			},
			want: []TaxInformation{
				{
					TotalIncome: 1000000,
					WHT:         100000,
					Allowances: []Allowance{
						{
							Type:   AllowanceTypeDonation,
							Amount: 10000,
						},
					},
				},
				{
					TotalIncome: 2000000,
					WHT:         200000,
					Allowances: []Allowance{
						{
							Type:   AllowanceTypeDonation,
							Amount: 20000,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			cr := NewCSVReader(nil)
			got, err := cr.parseTaxRecords(tc.records)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestParseTaxRecords_Error(t *testing.T) {

	t.Run("invalid header", func(t *testing.T) {
		// Arrange
		records := [][]string{
			{"invalid", "header"},
			{"1000000", "100000", "10000"},
		}

		// Act
		cr := NewCSVReader(nil)
		_, err := cr.parseTaxRecords(records)

		// Assert
		assert.Error(t, err)
	})

	t.Run("mis-spelled header", func(t *testing.T) {
		// Arrange
		records := [][]string{
			{"totalIncoming", "wht", "donation"},
			{"1000000", "100000", "10000"},
		}

		// Act
		cr := NewCSVReader(nil)
		_, err := cr.parseTaxRecords(records)

		// Assert
		assert.Error(t, err)
	})

	t.Run("invalid row", func(t *testing.T) {
		// Arrange
		records := [][]string{
			{"totalIncome", "wht", "donation"},
			{"1000000", "100000"},
		}

		// Act
		cr := NewCSVReader(nil)
		_, err := cr.parseTaxRecords(records)

		// Assert
		assert.Error(t, err)
	})

	t.Run("totalIncome is not number", func(t *testing.T) {
		// Arrange
		records := [][]string{
			{"totalIncome", "wht", "donation"},
			{"string", "100000", "10000"},
		}

		// Act
		cr := NewCSVReader(nil)
		_, err := cr.parseTaxRecords(records)

		// Assert
		assert.Error(t, err)
	})

	t.Run("wht is not number", func(t *testing.T) {
		// Arrange
		records := [][]string{
			{"totalIncome", "wht", "donation"},
			{"1000000", "string", "10000"},
		}

		// Act
		cr := NewCSVReader(nil)
		_, err := cr.parseTaxRecords(records)

		// Assert
		assert.Error(t, err)
	})

	t.Run("donation is not number", func(t *testing.T) {
		// Arrange
		records := [][]string{
			{"totalIncome", "wht", "donation"},
			{"1000000", "100000", "string"},
		}

		// Act
		cr := NewCSVReader(nil)
		_, err := cr.parseTaxRecords(records)

		// Assert
		assert.Error(t, err)
	})
}

func TestReadRecords(t *testing.T) {
	// Arrange
	t.Run("success", func(t *testing.T) {
		data := "totalIncome,wht,donation" + "\n"
		data += "1000000,100000,10000" + "\n"
		data += "2000000,200000,20000"
		want := []TaxInformation{
			{
				TotalIncome: 1000000,
				WHT:         100000,
				Allowances: []Allowance{
					{
						Type:   AllowanceTypeDonation,
						Amount: 10000,
					},
				},
			},
			{
				TotalIncome: 2000000,
				WHT:         200000,
				Allowances: []Allowance{
					{
						Type:   AllowanceTypeDonation,
						Amount: 20000,
					},
				},
			},
		}
		r := bytes.NewReader([]byte(data))

		// Act
		cr := NewCSVReader(r)
		got, err := cr.readRecords()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("error", func(t *testing.T) {
		// Arrange
		data := "totalIncome,wht,donation" + "\n"
		data += "1000000,100000" + "\n"
		r := bytes.NewReader([]byte(data))

		// Act
		cr := NewCSVReader(r)
		_, err := cr.readRecords()

		// Assert
		assert.Error(t, err)
	})
}
