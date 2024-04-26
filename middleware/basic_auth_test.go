//go:build unit

package middleware

import (
	"encoding/base64"
	"github.com/golfz/assessment-tax/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	// Arrange
	testCases := []struct {
		name     string
		username string
		password string
		want     int
	}{
		{
			name:     "correct username and password",
			username: "admin",
			password: "correct",
			want:     http.StatusOK,
		},
		{
			name:     "wrong username and password",
			username: "admin",
			password: "wrong",
			want:     http.StatusUnauthorized,
		},
	}

	cfg := config.Config{
		AdminUsername: "admin",
		AdminPassword: "correct",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			basicAuthMiddleware := BasicAuth(cfg)
			e.Use(middleware.BasicAuth(basicAuthMiddleware))
			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "protected data")
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			authHeader := "basic " + base64.StdEncoding.EncodeToString([]byte(tc.username+":"+tc.password))
			req.Header.Set(echo.HeaderAuthorization, authHeader)
			rec := httptest.NewRecorder()

			// Act
			e.ServeHTTP(rec, req)

			// Assert
			assert.Equal(t, tc.want, rec.Code)
			if rec.Code == http.StatusOK {
				assert.Equal(t, "protected data", rec.Body.String())
			}
		})
	}
}
