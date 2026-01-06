package api

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode) // set gin to test mode
	// the reason is that in debug mode, gin will print many logs in console which is not what we want in tests
	os.Exit(m.Run())
}