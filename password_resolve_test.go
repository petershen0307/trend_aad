package trendaad

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_password_fromEnv(t *testing.T) {
	os.Setenv("TREND_PASSWORD", "ps")
	defer os.Unsetenv("TREND_PASSWORD")
	password := retrievePassword([]string{})
	assert.Equal(t, "ps", password)
}

func Test_password_fromArgs(t *testing.T) {
	password := retrievePassword([]string{"program name", "-p", "ps"})
	assert.Equal(t, "ps", password)
}
