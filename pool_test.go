package chromedp

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
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

	pool := FromContext(taskCtx).pool.(*ExecPool)
	tempDir := pool.flags["user-data-dir"].(string)

	cancel()
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond)
		if _, err := os.Lstat(tempDir); os.IsNotExist(err) {
			return
		}
	}
	t.Fatalf("temporary user data dir %q not deleted", tempDir)
}
