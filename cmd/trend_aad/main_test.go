package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_userName_fromEnv(t *testing.T) {
	os.Setenv("TREND_USERNAME", "test")
	defer os.Unsetenv("TREND_USERNAME")
	user := retrieveUser([]string{})
	assert.Equal(t, "test", user)
}

func Test_userName_fromArgs(t *testing.T) {
	user := retrieveUser([]string{"program name", "test"})
	assert.Equal(t, "test", user)
}

func Test_userName_fromStdIn(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdin = r
	go func() {
		fmt.Fprintln(w, "test")
		w.Close()
	}()
	user := retrieveUser([]string{})
	assert.Equal(t, "test", user)
}

func Test_password_fromEnv(t *testing.T) {
	os.Setenv("TREND_PASSWORD", "ps")
	defer os.Unsetenv("TREND_PASSWORD")
	password := retrievePassword([]string{})
	assert.Equal(t, "ps", password)
}

func Test_password_fromArgs(t *testing.T) {
	password := retrievePassword([]string{"program name", "user", "ps"})
	assert.Equal(t, "ps", password)
}