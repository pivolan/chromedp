package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	dockerPoolOpts := []chromedp.DockerPoolOption{}

	ctxt, cancel := chromedp.NewPool(context.Background(), chromedp.WithDockerPool(dockerPoolOpts...))
	defer cancel()

	task1Context, cancel := chromedp.NewContext(ctxt)
	defer cancel()

	if err := chromedp.Run(task1Context, myTask()); err != nil {
		log.Fatal(err)
	}
}

func myTask() chromedp.Tasks {
	return []chromedp.Action{}
}
