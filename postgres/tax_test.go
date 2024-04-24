package postgres

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golfz/assessment-tax/tax"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDeduction_Success(t *testing.T) {
	testcases := []struct {
		name string
		rows *sqlmock.Rows
		want tax.Deduction
	}{
		{
			name: "found all deductions, expect deduction according to rows",
			rows: sqlmock.NewRows([]string{"name", "amount"}).
				AddRow("personal", "60000.00").
				AddRow("k-receipt", "50000.00").
				AddRow("donation", "100000.00"),
			want: tax.Deduction{
				Personal: 60000.00,
				KReceipt: 50000.00,
				Donation: 100000.00,
			},
		},
		{
			name: "not found any deductions, expect deduction with zero value",
			rows: sqlmock.NewRows([]string{"name", "amount"}),
			want: tax.Deduction{
				Personal: 0.0,
				KReceipt: 0.0,
				Donation: 0.0,
			},
		},
		{
			name: "found unknown deductions, expect deduction with zero value",
			rows: sqlmock.NewRows([]string{"name", "amount"}).
				AddRow("unknown", "10000.00"),
			want: tax.Deduction{
				Personal: 0.0,
				KReceipt: 0.0,
				Donation: 0.0,
			},
		},
		{
			name: "found some deductions, expect deduction according to row found and other with zero value",
			rows: sqlmock.NewRows([]string{"name", "amount"}).
				AddRow("k-receipt", "50000.00"),
			want: tax.Deduction{
				Personal: 0.0,
				KReceipt: 50000.00,
				Donation: 0.0,
			},
		},
	}

	for _, tc := range testcases {
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
			deduction, err := pg.GetDeduction()

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.want, deduction)
		})
	}
}

func TestGetDeduction_Error(t *testing.T) {
	testcases := []struct {
		name      string
		rows      *sqlmock.Rows
		expectSQL string
		want      tax.Deduction
	}{
		{
			name: "cannot query from unknown_table expect error",
			rows: sqlmock.NewRows([]string{"name", "amount"}).
				AddRow("personal", "60000.00").
				AddRow("k-receipt", "50000.00").
				AddRow("donation", "100000.00").
				AddRow("ignore_field", "10000.00"),
			expectSQL: `SELECT name, amount FROM unknown_table`,
			want:      tax.Deduction{},
		},
		{
			name: "cannot scan row because not enough column expect error",
			rows: sqlmock.NewRows([]string{"name"}).
				AddRow("personal").
				AddRow("k-receipt").
				AddRow("donation"),
			expectSQL: `SELECT name, amount FROM deductions`,
			want:      tax.Deduction{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			mock.ExpectQuery(tc.expectSQL).WillReturnRows(tc.rows)
			pg := Postgres{Db: db}

			// Act
			deduction, err := pg.GetDeduction()

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tc.want, deduction)
		})
	}
}
