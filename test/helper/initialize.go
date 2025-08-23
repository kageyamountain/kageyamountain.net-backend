package helper

import (
	"sync"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var once sync.Once

func InitializeIntegrationTest(t testing.TB) {
	t.Helper()

	once.Do(func() {
		loadEnvFile(t)
	})
}

func loadEnvFile(t testing.TB) {
	t.Helper()

	err := godotenv.Load("../../.env.dev")
	require.NoError(t, err)
}
