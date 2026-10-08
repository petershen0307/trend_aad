package trendaad

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/sirupsen/logrus"
)

func InitialBrowser() *rod.Browser {
	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Panicf("Can't get user home dir: %v", err)
	}
	// ponytail: single shared profile, per-user profiles if multi-account matters
	dir := filepath.Join(home, ".trend_aad", "chrome-profile")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		logrus.Panicf("Can't create browser profile dir(%s): %v", dir, err)
	}
	l := launcher.New().Leakless(false).UserDataDir(dir)
	exec.Command("pkill", "-f", dir).Run()
	os.Remove(filepath.Join(dir, "SingletonLock"))
	os.Remove(filepath.Join(dir, "SingletonSocket"))
	if path, exists := launcher.LookPath(); exists {
		logrus.Infof("detect browser: %s", path)
		l = l.Bin(path)
	} else {
		// try to install chromium
		logrus.Info("install browser")
	}
	controlURL := l.MustLaunch()
	browser := rod.New().ControlURL(controlURL).MustConnect().Logger(logrus.StandardLogger())
	logrus.Debugf("server monitor url: %s", browser.ServeMonitor(""))
	return browser
}
