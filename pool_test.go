package chromedp

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestExecPool(t *testing.T) {
	t.Parallel()

	poolCtx, cancel := NewPool(context.Background(), WithExecPool(poolOpts...))
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

	tempDir := FromContext(taskCtx).browser.UserDataDir
	pool := FromContext(taskCtx).Pool

	cancel()
	pool.Wait()

	if _, err := os.Lstat(tempDir); os.IsNotExist(err) {
		return
	}
	t.Fatalf("temporary user data dir %q not deleted", tempDir)
}
