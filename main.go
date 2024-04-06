package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"golang.org/x/term"
)

func main() {
	controlURL := ""
	if path, exists := launcher.LookPath(); exists {
		log.Println("detect browser", path)
		controlURL = launcher.New().Bin(path).Headless(true).Leakless(false).MustLaunch()
	} else {
		// try to install chromium
		controlURL = launcher.New().Headless(true).MustLaunch()
	}
	browser := rod.New().ControlURL(controlURL).MustConnect()
	launcher.Open(browser.ServeMonitor(""))
	defer browser.MustClose()

	trendAADURL := "https://awssts.infosec.trendmicro.com"
	user := os.Args[1]
	awsAccount := os.Args[2]
	fmt.Print("Enter Password: ")
	// to hide user input on terminal
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("")
	password := string(bytePassword)
	// fmt.Println(aadURL, user, password)

	// it will receive a 302 redirection
	page := browser.MustPage(trendAADURL).MustWaitStable()

	// cursor to user account then click
	page.MustElement("#i0116").MustInput(user).MustType(input.Enter)
	// page.MustElement("#idSIButton9").MustClick()

	// cursor to password then click
	page.MustWaitStable().MustElement("#i0118").MustInput(password).MustType(input.Enter)
	// page.MustWaitStable().MustElement("#idSIButton9").MustClick()

	// show number for authenticator
	number := page.MustWaitStable().MustElement("#idRichContext_DisplaySign").MustText()
	fmt.Println("authenticator", number)

	// wait page redirect
	for !strings.Contains(page.MustInfo().URL, trendAADURL) {
		// log.Println("wait url redirect")
		time.Sleep(1 * time.Second)
	}

	// expand aws account div
	page.MustElement(fmt.Sprintf("#accordion%s > div > a", awsAccount)).MustClick()
	// select admin
	page.MustElement(fmt.Sprintf("#collapse%s > div > p > button", awsAccount)).MustClick()
	// get sts json format
	page.MustWaitStable().MustElement("#json-tab").MustClick()
	text := page.MustWaitStable().MustElement("#copyTextjson").MustClick().MustText()
	fmt.Println(text)
	aksk := AKSK{}
	if err := json.Unmarshal([]byte(text), &aksk); err != nil {
		log.Panicln("json unmarshal error=", err)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln(err)
	}
	const awsDir = ".aws"
	const awsCredentialPath = awsDir + "/" + "credentials"
	awsCredential := fmt.Sprintf(`[default]
aws_access_key_id = %s
aws_secret_access_key = %s
aws_session_token  = %s
`, aksk.AccessKeyId, aksk.SecretAccessKey, aksk.SessionToken)
	if _, err := os.Stat(filepath.Join(homeDir, awsDir)); errors.Is(err, os.ErrNotExist) {
		log.Println(err)
		if err := os.MkdirAll(filepath.Join(homeDir, awsDir), 0777); err != nil {
			log.Panicln("create folder error=", err)
		}
	}
	if err := os.WriteFile(filepath.Join(homeDir, awsCredentialPath), []byte(awsCredential), 0666); err != nil {
		log.Panicln("write file error=", err)
	}
}

type AKSK struct {
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
}
