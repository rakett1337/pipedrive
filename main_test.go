package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("TEST_MODE", "true")

	exitCode := m.Run()
	os.Exit(exitCode)
}
