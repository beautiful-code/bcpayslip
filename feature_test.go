package main

import (
	"context"
	"log"
	"testing"

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
		cdp.Navigate(`http://localhost:3001`),
		cdp.WaitVisible("#google-signin", cdp.ByID),
	})
	if err != nil {
		log.Fatal(err)
	}
	tearDown(c, ctxt)
}

/* func TestHomePage(t *testing.T) {
 *     c, ctxt, cancel := setup()
 *     defer cancel()
 *     login := "#google-signin"
 *     err := c.Run(ctxt, cdp.Tasks{
 *         cdp.Navigate(`http://localhost:3001`),
 *         cdp.Click(login),
 *         cdp.Sleep(5 * time.Second),
 *         cdp.WaitVisible("#submit", cdp.ByID),
 *     })
 *     if err != nil {
 *         log.Fatal(err)
 *     }
 *     tearDown(c, ctxt)
 * } */
