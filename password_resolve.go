package trendaad

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

func retrievePassword(args []string) string {
	password := os.Getenv("TREND_PASSWORD")
	if len(args) >= 3 {
		for i, arg := range args {
			if arg == "-p" && i+1 <= len(args) {
				password = args[i+1]
				break
			}
		}
	}
	if password == "" {
		fmt.Print("Enter Password: ")
		// to hide user input on terminal
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			logrus.Error(err)
			return ""
		}
		fmt.Println("")
		password = string(bytePassword)
	}
	return password
}
