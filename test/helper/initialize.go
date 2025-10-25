package helper

import (
	"os"
	"sync"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var once sync.Once

func InitializeIntegrationTest(tb testing.TB) {
	tb.Helper()

	once.Do(func() {
		loadEnvFile(tb)
	})
}

func loadEnvFile(tb testing.TB) {
	tb.Helper()

	env := os.Getenv("ENV")

	if env == "ci" {
		// CI環境
		err := godotenv.Load("../../.env.ci")
		require.NoError(tb, err)
	} else {
		// ローカル開発環境
		err := godotenv.Load("../../.env.dev")
		require.NoError(tb, err)
	}
}
