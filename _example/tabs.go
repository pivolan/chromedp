package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	ctxt, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// create new tab
	tabContext, _ := chromedp.NewContext(ctxt)

	// runs in first tab
	err := chromedp.Run(ctxt, myTask())
	if err != nil {
		log.Fatal(err)
	}

	// runs in second tab
	err = chromedp.Run(tabContext, myTask())
	if err != nil {
		log.Fatal(err)
	}
}

func myTask() chromedp.Tasks {
	return []chromedp.Action{}
}
