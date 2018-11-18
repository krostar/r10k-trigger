package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	var cfg *appConfig

	require.NotPanics(t, func() {
		os.Setenv("APP_HTTPSERVER_LISTENADDRESS", ":9988") // nolint: errcheck, gosec
		cfg = initConfig("../../configs/trigger-api/default.yaml")
	})

	require.NotNil(t, cfg)
	require.Equal(t, ":9988", cfg.HTTPServer.ListenAddress)
}

func TestInitConfig_fail(t *testing.T) {
	require.Panics(t, func() {
		os.Setenv("APP_GRACEFULSHUTDOWNTIMEOUT", "[42]") // nolint: errcheck, gosec
		initConfig("../../configs/trigger-api/default.yaml")
	})
}
