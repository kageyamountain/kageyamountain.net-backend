package helper

import (
	"os"
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

	env := os.Getenv("ENV")

	if env == "ci" {
		// CI環境
		err := godotenv.Load("../../.env.ci")
		require.NoError(t, err)
	} else {
		// ローカル開発環境
		err := godotenv.Load("../../.env.dev")
		require.NoError(t, err)
	}
}
