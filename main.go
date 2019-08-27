//get debt information from to sites

package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

//site information
type WebAccount struct {
	Url     string
	Login   string
	SelName string
	Passwd  string
	SelPass string
}

var kvado = WebAccount{
	Url:     "https://cabinet.kvado.ru/login",
	Login:   "omeen@bk.ru",
	SelName: `//input[@id="LoginForm_email"]`,
	Passwd:  "password123456",
	SelPass: `//input[@id="LoginForm_password"]`,
}

var rksEnergo = WebAccount{
	Url:     "http://lk.rks-energo.ru",
	Login:   "omeen@bk.ru",
	SelName: `document.querySelector("#P101_USERNAME")`,
	Passwd:  "password123456",
	SelPass: `document.querySelector("#P101_PASSWORD")`,
}

//get information from sites
func main() {
	rksEnergoGetDebt()
	time.Sleep(4 * time.Second)
	kvadoGetDebt()
	time.Sleep(10 * time.Second) // TODO: check opening main page
}

//get debt for rks-energo.ru
func rksEnergoGetDebt() {
	log.Println("start scraping")
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false), // browser will be visible
		//chromedp.ProxyServer("http://192.168.55.53:3128/"),
		//chromedp.UserAgent(`Mozilla/5.0 (iPhone; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.25 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1`),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err := chromedp.Run(ctx, WebLoginJs(rksEnergo))
	if err != nil {
		log.Fatal(err)
	}
	err = chromedp.Run(ctx, ClickCharges())
	if err != nil {
		log.Println("error=", err)
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	err = chromedp.Run(ctx, GetInnerText(`document.querySelector("#apexir_DATA_PANEL").innerText`))
	if err != nil {
		log.Fatal(err)
	}
}

//get debt from cabinet.kvado.ru
func kvadoGetDebt() {
	log.Println("start scraping")
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", true), // browser will be visible
		//chromedp.ProxyServer("http://192.168.55.53:3128/"),
		chromedp.UserAgent("Opera"))
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := chromedp.Run(ctx, WebLogin(kvado))
	if err != nil {
		log.Fatal(err)
	}
	err = chromedp.Run(ctx, ReadSelector("yw0"))
	if err != nil {
		log.Fatal(err)
	}
	err = chromedp.Run(ctx, ReadSelector(`#navbar-container > div.navbar-buttons.navbar-header.pull-right > ul > li.balance-container.red > a`))
	if err != nil {
		log.Fatal(err)
	}

}

//get text from titlename print to log
func ReadTitle(titlename string) chromedp.Tasks {
	var res string
	return chromedp.Tasks{
		chromedp.WaitVisible(titlename),
		chromedp.Text(titlename, &res, chromedp.NodeVisible, chromedp.ByID),
		chromedp.ActionFunc(func(context.Context) error {
			log.Printf(res)
			return nil
		}),
	}
}

//login to site using pressing enter after writing password
func WebLoginJs(a WebAccount) chromedp.Tasks {
	log.Println("WebLogin start")
	return chromedp.Tasks{
		chromedp.Navigate(a.Url),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitVisible(a.SelName, chromedp.ByJSPath),
		chromedp.SendKeys(a.SelName, a.Login, chromedp.ByJSPath),
		//send keys with Enter
		chromedp.SendKeys(a.SelPass, a.Passwd+"\r\n", chromedp.ByJSPath),
		chromedp.ActionFunc(func(context.Context) error {
			log.Printf("subbmiting login form")
			return nil
		}),
	}
}

//submitting login form
func WebLogin(a WebAccount) chromedp.Tasks {
	log.Println("WebLogin start")
	return chromedp.Tasks{
		chromedp.Navigate(a.Url),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(a.SelPass),
		chromedp.SendKeys(a.SelName, a.Login),
		chromedp.SendKeys(a.SelPass, a.Passwd),
		chromedp.Submit(a.SelPass),
		chromedp.ActionFunc(func(context.Context) error {
			log.Printf("subbmiting login form")
			return nil
		}),
	}
}

//print to log text
func ReadSelector(selector string) chromedp.Tasks {
	var res string
	return chromedp.Tasks{
		chromedp.WaitVisible(selector),
		chromedp.Text(selector, &res, chromedp.NodeVisible, chromedp.ByID),
		chromedp.ActionFunc(func(context.Context) error {
			log.Printf(res)
			return nil
		}),
	}
}

//click to button начисления
func ClickCharges() chromedp.Tasks {
	log.Println("clickCharges")
	var res bool
	return chromedp.Tasks{
		chromedp.Evaluate(`javascript:apex.submit('Начисления');`, &res),
		chromedp.ActionFunc(func(context.Context) error {
			log.Println("result executing clickCharges", res)
			return nil
		}),
	}
}

//get text using java script eval function
func GetInnerText(jsSelector string) chromedp.Tasks {
	log.Println("GetInnerText")
	var res interface{}
	return chromedp.Tasks{
		chromedp.Evaluate(jsSelector, &res),
		chromedp.ActionFunc(func(context.Context) error {
			log.Println("res=", res)
			return nil
		}),
	}
}
