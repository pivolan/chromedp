package chromedp

import (
	"context"
	"log"
	"testing"
)

func TestExecPool(t *testing.T) {
	t.Parallel()

	poolCtx, cancel := NewPool(context.Background(), WithExecPool())
	defer cancel()

	// TODO: test that multiple child contexts are run in different
	// processes and browsers.

	taskCtx, cancel := NewContext(poolCtx, WithURL(testdataDir+"/form.html"))
	defer cancel()

	want := "insert"
	var got string
	if err := Run(taskCtx, Text("#foo", &got, ByID)); err != nil {
		log.Fatal(err)
	}
	if got != want {
		log.Fatalf("wanted %q, got %q", want, got)
	}
}
