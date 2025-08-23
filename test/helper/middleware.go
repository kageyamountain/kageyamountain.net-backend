package helper

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// テスト用にGinのコンテキストに値を設定するミドルウェア
// ミドルウェアでgin.Contextにセットされる値を、それらのミドルウェアを呼ぶことなくセットするために利用する
func InjectGinContext(t *testing.T, values map[string]any) gin.HandlerFunc {
	t.Helper()

	return func(c *gin.Context) {
		for k, v := range values {
			c.Set(k, v)
		}
		c.Next()
	}
}
