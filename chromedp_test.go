package chromedp

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"
)

var (
	testdataDir string

	poolCtx context.Context

	poolOpts = []ExecAllocatorOption{
		NoFirstRun,
		NoDefaultBrowserCheck,
		Headless,
	}
)

func testAllocate(t *testing.T, path string) (_ context.Context, cancel func()) {
	ctx, cancel := NewContext(poolCtx, WithURL(testdataDir+"/"+path))

	//if err := WithLogf(t.Logf)(c.c); err != nil {
	//        t.Fatalf("could not set logf: %v", err)
	//}
	//if err := WithDebugf(t.Logf)(c.c); err != nil {
	//        t.Fatalf("could not set debugf: %v", err)
	//}
	//if err := WithErrorf(t.Errorf)(c.c); err != nil {
	//        t.Fatalf("could not set errorf: %v", err)
	//}

	return ctx, cancel
}

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("could not get working directory: %v", err))
	}
	testdataDir = "file://" + path.Join(wd, "testdata")

	// it's worth noting that newer versions of chrome (64+) run much faster
	// than older ones -- same for headless_shell ...
	if execPath := os.Getenv("CHROMEDP_TEST_RUNNER"); execPath != "" {
		poolOpts = append(poolOpts, ExecPath(execPath))
	}
	// not explicitly needed to be set, as this vastly speeds up unit tests
	if noSandbox := os.Getenv("CHROMEDP_NO_SANDBOX"); noSandbox != "false" {
		poolOpts = append(poolOpts, NoSandbox)
	}
	// must be explicitly set, as disabling gpu slows unit tests
	if disableGPU := os.Getenv("CHROMEDP_DISABLE_GPU"); disableGPU != "" && disableGPU != "false" {
		poolOpts = append(poolOpts, DisableGPU)
	}

	ctx, cancel := NewAllocator(context.Background(), WithExecAllocator(poolOpts...))
	poolCtx = ctx

	code := m.Run()

	cancel()
	FromContext(ctx).Wait()
	os.Exit(code)
}
