package trendaad

import (
	"runtime"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/sirupsen/logrus"
)

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
