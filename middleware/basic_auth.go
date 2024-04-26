package middleware

import (
	"crypto/subtle"
	"github.com/golfz/assessment-tax/config"
	"github.com/labstack/echo/v4"
)

func BasicAuth(cfg config.Config) func(username, password string, c echo.Context) (bool, error) {
	return func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.AdminUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(cfg.AdminPassword)) == 1 {
			return true, nil
		}
		return false, nil
	}
}
