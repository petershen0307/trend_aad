package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"golang.org/x/term"
)

func main() {
	l := launcher.New().Headless(false)
	defer l.Cleanup()
	url := l.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	launcher.Open(browser.ServeMonitor(""))
	defer browser.MustClose()

	aadURL := os.Args[1]
	user := os.Args[2]
	fmt.Print("Enter Password: ")
	// to hide user input on terminal
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println(err)
		return
	}
	password := string(bytePassword)
	// fmt.Println(aadURL, user, password)

	// it will receive a 302 redirection
	page := browser.MustPage(aadURL).MustWaitStable()

	// cursor to user account then click
	page.MustElement("#i0116").MustInput(user).MustType(input.Enter)
	// page.MustElement("#idSIButton9").MustClick()

	// cursor to password then click
	page.MustWaitStable().MustElement("#i0118").MustInput(password).MustType(input.Enter)
	// page.MustWaitStable().MustElement("#idSIButton9").MustClick()

	// wait page redirect
	for !strings.Contains(page.MustInfo().URL, aadURL) {
		// log.Println("wait url redirect")
		time.Sleep(1 * time.Second)
	}
	// expand dspm int
	page.MustElement("#accordion530748458907 > div > a").MustClick()
	// select admin
	page.MustElement("#collapse530748458907 > div > p > button").MustClick()
	// get sts
	text := page.MustWaitStable().MustElement("#copyTextbash").MustClick().MustText()
	fmt.Println(strings.ReplaceAll(strings.ReplaceAll(text, "export ", ""), ";", ""))
}
