package trendaad

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/sirupsen/logrus"
)

const trendAwsStsUrl = "https://awssts.infosec.trendmicro.com"

func InitialBrowser() *rod.Browser {
	controlURL := ""
	if path, exists := launcher.LookPath(); exists {
		logrus.Info("detect browser", path)
		controlURL = launcher.New().Bin(path).Headless(true).MustLaunch()
	} else {
		// try to install chromium
		logrus.Info("install browser")
		controlURL = launcher.New().Headless(true).MustLaunch()
	}
	browser := rod.New().ControlURL(controlURL).MustConnect().Logger(logrus.StandardLogger())
	launcher.Open(browser.ServeMonitor(""))
	runtime.SetFinalizer(browser, func(browser *rod.Browser) {
		browser.MustClose()
	})
	return browser
}

func LoginPage(browser *rod.Browser, user, password string) *rod.Page {
	logrus.Info("Wait for the login page")
	// it will receive a 302 redirection
	page := browser.MustPage(trendAwsStsUrl).MustWaitStable()

	// cursor to user account then click
	page.MustElement("#i0116").MustInput(user).MustType(input.Enter)

	// cursor to password then click
	page.MustWaitStable().MustElement("#i0118").MustInput(password).MustType(input.Enter)

	// show number for authenticator
	logrus.Info("Wait for the two factor authenticator number")
	number := page.MustWaitStable().MustElement("#idRichContext_DisplaySign").MustText()
	fmt.Println("authenticator", number)

	// wait page redirect
	for !strings.Contains(page.MustInfo().URL, trendAwsStsUrl) {
		time.Sleep(100 * time.Millisecond)
	}
	return page
}

func ExtractAwsStsFromPage(page *rod.Page) TrendAwsSts {
	// get credential from resource F12>Application>Frames>Top>https://awssts.infosec.trendmicro.com
	pageHtml, err := page.GetResource(trendAwsStsUrl)
	if err != nil {
		logrus.Errorf("Can't get the resource(%s): %v", trendAwsStsUrl, err)
	}
	scanner := bufio.NewScanner(strings.NewReader(html.UnescapeString(string(pageHtml))))
	creds := TrendAwsSts{}
	// the credential in page is a json object start with this prefix in a line
	const prefix = "var creds = "
	for scanner.Scan() {
		if s := strings.TrimSpace(scanner.Text()); strings.HasPrefix(s, prefix) {
			s = strings.TrimPrefix(s, prefix)
			s = strings.TrimRight(s, ";")
			logrus.Debug(s)
			err := json.Unmarshal([]byte(s), &creds)
			if err != nil {
				logrus.Error(err)
			}
			break
		}
	}
	return creds
}

/*
example of TrendAwsSts:

	{
		"123456789012": {
			"AAD-READONLY_123456789012: {
				"AccessKeyId": "",
				"SecretAccessKey": "",
				"SessionToken": "",
				"Expiration": "2025-01-02T06:53:16+00:00"
			},
			"name": "test"
		}
	}
*/
type TrendAwsSts map[string]map[string]interface{}

type AwsSts struct {
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Expiration      string `json:"Expiration"`
}

func (trendAwsSts TrendAwsSts) FlushAwsCredential(awsCredentialFile io.StringWriter) {
	for _, accountVals := range trendAwsSts {
		if _, ok := accountVals["name"]; !ok {
			continue
		}
		accountName := accountVals["name"].(string)
		for k, v := range accountVals {
			if k == "name" {
				continue
			}
			byteData, err := json.Marshal(v)
			if err != nil {
				logrus.Errorf("Can't marshal the credential: %v", err)
			}
			var sts AwsSts
			err = json.Unmarshal(byteData, &sts)
			if err != nil {
				logrus.Debugf("Skip data(%s): %v", string(byteData), err)
				continue
			}
			iniFormat := convertToIniFormat(fmt.Sprintf("%s_%s", accountName, k), sts)
			_, err = awsCredentialFile.WriteString(iniFormat + "\n")
			if err != nil {
				logrus.Errorf("Can't write the credential file: %v", err)
			}
		}
	}
}

func convertToIniFormat(section string, sts AwsSts) string {
	const template = `[%s]
aws_access_key_id = %s
aws_secret_access_key = %s
aws_session_token = %s`
	return fmt.Sprintf(template,
		section, sts.AccessKeyId, sts.SecretAccessKey, sts.SessionToken)
}
