package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	trendaad "github.com/petershen0307/trend_aad"
	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
)

func main() {
	initLogger()
	browser := trendaad.InitialBrowser()
	defer browser.MustClose()
	var page *rod.Page
	if isSemiAuto(os.Args) {
		page = trendaad.WaitSemiAutoLogin(browser)
	} else {
		user := retrieveUser(os.Args)
		page = trendaad.LoginPage(browser, user)
	}
	sts := trendaad.ExtractAwsStsFromPage(page)
	awsCredentialFile, err := openAwsCredentialFile()
	if err != nil {
		return
	}
	defer awsCredentialFile.Close()
	logrus.Debugf("aws credential file: %v", awsCredentialFile.Name())
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

func isSemiAuto(args []string) bool {
	return slices.Contains(args, "-rod=show")
}

func retrieveUser(args []string) string {
	user := os.Getenv("TREND_USERNAME")
	if len(args) >= 3 {
		for i, arg := range args {
			if arg == "-u" && i+1 <= len(args) {
				user = args[i+1]
				break
			}
		}
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
