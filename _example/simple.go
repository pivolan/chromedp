package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithURL("https://github.com/"))
	defer cancel()

	if err := chromedp.Run(ctx, myTask()); err != nil {
		log.Fatal(err)
	}
	// TODO: make this unnecessary
	cancel()
	time.Sleep(time.Millisecond)
}

func myTask() chromedp.Tasks {
	return []chromedp.Action{
		chromedp.Sleep(10 * time.Second),
	}
}
