package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	trendaad "github.com/petershen0307/trend_aad"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

func main() {
	initLogger()
	user := retrieveUser(os.Args)
	password := retrievePassword(os.Args)
	browser := trendaad.InitialBrowser()
	page := trendaad.LoginPage(browser, user, password)
	sts := trendaad.ExtractAwsStsFromPage(page)
	awsCredentialFile, err := openAwsCredentialFile()
	if err != nil {
		return
	}
	defer awsCredentialFile.Close()
	sts.FlushAwsCredential(awsCredentialFile)
}

func initLogger() {
	level := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(level) {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func retrieveUser(args []string) string {
	user := os.Getenv("TREND_USERNAME")
	if len(args) >= 2 {
		user = args[1]
	}
	// let user input the user name
	if user == "" {
		fmt.Print("Enter User: ")
		if _, err := fmt.Scanln(&user); err != nil {
			logrus.Error(err)
		}
	}
	return user
}

func retrievePassword(args []string) string {
	password := os.Getenv("TREND_PASSWORD")
	if len(args) >= 3 {
		password = args[2]
	}
	if password == "" {
		fmt.Print("Enter Password: ")
		// to hide user input on terminal
		bytePassword, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			logrus.Error(err)
			return ""
		}
		fmt.Println("")
		password = string(bytePassword)
	}
	return password
}

func openAwsCredentialFile() (*os.File, error) {
	// check aws credential file
	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	awsCredentialPath := filepath.Join(home, ".aws", "credentials")
	awsCredentialFile, err := os.OpenFile(awsCredentialPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("Can't open the credential file(%s): %v", awsCredentialPath, err)
		return nil, err
	}
	return awsCredentialFile, nil
}
