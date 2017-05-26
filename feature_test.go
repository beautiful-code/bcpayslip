package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
	cdp "github.com/knq/chromedp"
)

func setup() (*cdp.CDP, context.Context, context.CancelFunc) {
	var err error
	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	// create chrome instance
	c, err := cdp.New(ctxt, cdp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}
	return c, ctxt, cancel
}

func tearDown(c *cdp.CDP, ctxt context.Context) {
	// shutdown chrome
	err := c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}
	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func TestLoginPage(t *testing.T) {
	c, ctxt, cancel := setup()
	defer cancel()
	err := c.Run(ctxt, cdp.Tasks{
		cdp.Navigate(os.Getenv("bc_test_host")),
		cdp.WaitVisible("#google-signin", cdp.ByID),
	})
	if err != nil {
		log.Fatal(err)
	}
	tearDown(c, ctxt)
}

func TestHomePage(t *testing.T) {
	c, ctxt, cancel := setup()
	defer cancel()
	login := "#google-signin"
	err := c.Run(ctxt, cdp.Tasks{
		cdp.Navigate(os.Getenv("bc_test_host")),
		cdp.Click(login),
		cdp.Sleep(2 * time.Second),
		cdp.WaitVisible("#identifierNext", cdp.ByID),
	})
	if err != nil {
		log.Fatal(err)
	}
	tearDown(c, ctxt)
}
