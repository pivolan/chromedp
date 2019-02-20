package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	ctxt, cancel := chromedp.NewContext(context.Background(), chromedp.WithURL("https://github.com/"))
	defer cancel()

	if err := chromedp.Run(ctxt, myTask()); err != nil {
		log.Fatal(err)
	}
}

func myTask() chromedp.Tasks {
	return []chromedp.Action{}
}
