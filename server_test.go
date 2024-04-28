//go:build unit

package main

import (
	"context"
	"github.com/golfz/assessment-tax/config"
	"github.com/labstack/echo/v4"
	"testing"
)

func TestStartServer(t *testing.T) {
	testCases := []struct {
		name string
		port int
		err  error
	}{
		{
			name: "valid port",
			port: 8080,
			err:  nil,
		},
		//{
		//	name: "invalid port",
		//	port: 0,
		//	err:  http.ErrServerClosed,
		//},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtx := context.Background()
			mockEcho := echo.New()
			mockEcho.HideBanner = true
			mockCfg := &config.Config{Port: tc.port}

			go startServer(mockEcho, mockCfg)

			err := mockEcho.Shutdown(mockCtx)
			if err != tc.err {
				t.Errorf("startServer() error = %v, wantErr %v", err, tc.err)
			}
		})
	}
}
