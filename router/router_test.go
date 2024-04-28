//go:build unit

package router

import (
	"database/sql"
	"github.com/golfz/assessment-tax/config"
	"github.com/golfz/assessment-tax/postgres"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	// Arrange
	e := New(&postgres.Postgres{DB: &sql.DB{}}, &config.Config{})

	req := httptest.NewRequest(http.MethodGet, "/not/registered/uri", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)
	r := e.Routes()

	// Assert
	assert.Greater(t, len(r), 0)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
