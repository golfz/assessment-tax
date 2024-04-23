package postgres

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golfz/assessment-tax/tax"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDeduction_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"name", "amount"}).
		AddRow("personal", "60000.00").
		AddRow("k-receipt", "50000.00").
		AddRow("donation", "100000.00").
		AddRow("unknown", "60000.00")
	mock.ExpectQuery(`SELECT name, amount FROM deductions`).WillReturnRows(rows)
	pg := Postgres{Db: db}
	want := tax.Deduction{
		Personal: 60000.00,
		KReceipt: 50000.00,
		Donation: 100000.00,
	}

	// Act
	deduction, err := pg.GetDeduction()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, want, deduction)
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
