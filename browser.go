package trendaad

import (
	"runtime"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/sirupsen/logrus"
)

func InitialBrowser() *rod.Browser {
	controlURL := ""
	l := launcher.New().Headless(true).Leakless(false)
	if path, exists := launcher.LookPath(); exists {
		logrus.Info("detect browser", path)
		controlURL = l.Bin(path).MustLaunch()
	} else {
		// try to install chromium
		logrus.Info("install browser")
		controlURL = l.MustLaunch()
	}
	browser := rod.New().ControlURL(controlURL).MustConnect().Logger(logrus.StandardLogger())
	runtime.SetFinalizer(browser, func(browser *rod.Browser) {
		browser.MustClose()
	})
	logrus.Debugf("server monitor url: %s", browser.ServeMonitor(""))
	return browser
}
