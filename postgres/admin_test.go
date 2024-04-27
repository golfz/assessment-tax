//go:build unit

package postgres

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetPersonalDeduction_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("^UPDATE (.+)").WillReturnResult(sqlmock.NewResult(0, 1))
	pg := Postgres{DB: db}

	// Act
	err = pg.SetPersonalDeduction(60000.00)

	// Assert
	assert.NoError(t, err)
	_ = mock.ExpectationsWereMet()
}

func TestSetPersonalDeduction_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("^UPDATE (.+)").WillReturnError(errors.New("unexpected error"))
	pg := Postgres{DB: db}

	// Act
	err = pg.SetPersonalDeduction(60000.00)

	// Assert
	assert.Error(t, err)
	_ = mock.ExpectationsWereMet()
}

func TestSetKReceiptDeduction_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("^UPDATE (.+)").WillReturnResult(sqlmock.NewResult(0, 1))
	pg := Postgres{DB: db}

	// Act
	err = pg.SetKReceiptDeduction(70000.00)

	// Assert
	assert.NoError(t, err)
	_ = mock.ExpectationsWereMet()
}
