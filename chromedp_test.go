package chromedp

import (
	"context"
	"log"
	"os"
	"path"
	"testing"
)

var (
	poolCtx context.Context

	testdataDir string

	//cliOpts = []runner.CommandLineOption{
	//        runner.NoDefaultBrowserCheck,
	//        runner.NoFirstRun,
	//}
)

func testAllocate(t *testing.T, path string) *Res {
	// TODO: cliOpts
	ctx, cancel := NewContext(poolCtx)
	_ = cancel // TODO

	if err := WithLogf(t.Logf)(c.c); err != nil {
		t.Fatalf("could not set logf: %v", err)
	}

	if err := WithDebugf(t.Logf)(c.c); err != nil {
		t.Fatalf("could not set debugf: %v", err)
	}

	if err := WithErrorf(t.Errorf)(c.c); err != nil {
		t.Fatalf("could not set errorf: %v", err)
	}

	//h := c.c.GetHandlerByIndex(0)
	//th, ok := h.(*TargetHandler)
	//if !ok {
	//        t.Fatalf("handler is invalid type")
	//}

	//th.logf, th.debugf = t.Logf, t.Logf
	//th.errf = func(s string, v ...interface{}) {
	//        t.Logf("TARGET HANDLER ERROR: "+s, v...)
	//}

	if path != "" {
		err = Run(ctx, Navigate(testdataDir+"/"+path))
		if err != nil {
			t.Fatalf("could not navigate to testdata/%s: %v", path, err)
		}
	}

	return c
}

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
		os.Exit(1)
	}
	testdataDir = "file://" + path.Join(wd, "testdata")

	/*
		// its worth noting that newer versions of chrome (64+) run much faster
		// than older ones -- same for headless_shell ...
		execPath := os.Getenv("CHROMEDP_TEST_RUNNER")
		if execPath == "" {
			execPath = runner.LookChromeNames("headless_shell")
		}
		cliOpts = append(cliOpts, runner.ExecPath(execPath))

		// not explicitly needed to be set, as this vastly speeds up unit tests
		if noSandbox := os.Getenv("CHROMEDP_NO_SANDBOX"); noSandbox != "false" {
			cliOpts = append(cliOpts, runner.NoSandbox)
		}
		// must be explicitly set, as disabling gpu slows unit tests
		if disableGPU := os.Getenv("CHROMEDP_DISABLE_GPU"); disableGPU != "" && disableGPU != "false" {
			cliOpts = append(cliOpts, runner.DisableGPU)
		}

		if targetTimeout := os.Getenv("CHROMEDP_TARGET_TIMEOUT"); targetTimeout != "" {
			defaultNewTargetTimeout, _ = time.ParseDuration(targetTimeout)
		}
		if defaultNewTargetTimeout == 0 {
			defaultNewTargetTimeout = 30 * time.Second
		}
	*/

	ctx, cancel := NewPool(context.Background())
	poolCtx = ctx

	code := m.Run()

	cancel()
	os.Exit(code)
}
