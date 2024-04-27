//go:build unit

package postgres

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golfz/assessment-tax/deduction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDeduction_Success(t *testing.T) {
	testCases := []struct {
		name string
		rows *sqlmock.Rows
		want deduction.Deduction
	}{
		{
			name: "found all deductions, expect deduction according to rows",
			rows: sqlmock.NewRows([]string{"name", "amount"}).
				AddRow("personal", "60000.00").
				AddRow("k-receipt", "50000.00").
				AddRow("donation", "100000.00"),
			want: deduction.Deduction{
				Personal: 60000.00,
				KReceipt: 50000.00,
				Donation: 100000.00,
			},
		},
		{
			name: "not found any deductions, expect deduction with zero value",
			rows: sqlmock.NewRows([]string{"name", "amount"}),
			want: deduction.Deduction{
				Personal: 0.0,
				KReceipt: 0.0,
				Donation: 0.0,
			},
		},
		{
			name: "found unknown deductions, expect deduction with zero value",
			rows: sqlmock.NewRows([]string{"name", "amount"}).
				AddRow("unknown", "10000.00"),
			want: deduction.Deduction{
				Personal: 0.0,
				KReceipt: 0.0,
				Donation: 0.0,
			},
		},
		{
			name: "found some deductions, expect deduction according to row found and other with zero value",
			rows: sqlmock.NewRows([]string{"name", "amount"}).
				AddRow("k-receipt", "50000.00"),
			want: deduction.Deduction{
				Personal: 0.0,
				KReceipt: 50000.00,
				Donation: 0.0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
			}
			defer db.Close()
			mock.ExpectQuery(`SELECT name, amount FROM deductions`).WillReturnRows(tc.rows)
			pg := Postgres{Db: db}

			// Act
			deductionData, err := pg.GetDeduction()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.want, deductionData)
		})
	}
}

func TestGetDeduction_Error(t *testing.T) {
	t.Run("query error, expect error with zero value", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectQuery(`SELECT name, amount FROM deductions`).WillReturnError(ErrCannotQueryDeduction)
		pg := Postgres{Db: db}
		wantDeduction := deduction.Deduction{}

		// Act
		gotDeduction, err := pg.GetDeduction()

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrCannotQueryDeduction)
		assert.Equal(t, wantDeduction, gotDeduction)
	})

	t.Run("scan row error (amount is not float), expect error with zero value", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		rows := sqlmock.NewRows([]string{"name", "amount"}).
			AddRow("personal", "abcdef").
			AddRow("k-receipt", "50000.00").
			AddRow("donation", "100000.00")
		mock.ExpectQuery(`SELECT name, amount FROM deductions`).WillReturnRows(rows)
		pg := Postgres{Db: db}
		wantDeduction := deduction.Deduction{}

		// Act
		gotDeduction, err := pg.GetDeduction()

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrCannotScanDeduction)
		assert.Equal(t, wantDeduction, gotDeduction)
	})
}
