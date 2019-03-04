package chromedp

import (
	"strings"
	"testing"
	"time"

	"github.com/chromedp/cdproto/page"
)

func TestNavigate(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	expurl, exptitle := testdataDir+"/image.html", "this is title"

	err = Run(ctx, Navigate(expurl))
	if err != nil {
		t.Fatal(err)
	}

	err = Run(ctx, WaitVisible(`#icon-brankas`, ByID))
	if err != nil {
		t.Fatal(err)
	}

	var urlstr string
	err = Run(ctx, Location(&urlstr))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(urlstr, expurl) {
		t.Errorf("expected to be on image.html, at: %s", urlstr)
	}

	var title string
	err = Run(ctx, Title(&title))
	if err != nil {
		t.Fatal(err)
	}
	if title != exptitle {
		t.Errorf("expected title to contain google, instead title is: %s", title)
	}
}

func TestNavigationEntries(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	tests := []string{
		"form.html",
		"image.html",
	}

	var entries []*page.NavigationEntry
	var index int64

	err = Run(ctx, NavigationEntries(&index, &entries))
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Errorf("expected to have 1 navigation entry: got %d", len(entries))
	}
	if index != 0 {
		t.Errorf("expected navigation index is 0, got: %d", index)
	}

	expIdx, expEntries := 1, 2
	for i, url := range tests {
		err = Run(ctx, Navigate(testdataDir+"/"+url))
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(50 * time.Millisecond)

		err = Run(ctx, NavigationEntries(&index, &entries))
		if err != nil {
			t.Fatal(err)
		}

		if len(entries) != expEntries {
			t.Errorf("test %d expected to have %d navigation entry: got %d", i, expEntries, len(entries))
		}
		if index != int64(i+1) {
			t.Errorf("test %d expected navigation index is %d, got: %d", i, i, index)
		}

		expIdx++
		expEntries++
	}
}

func TestNavigateToHistoryEntry(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	var entries []*page.NavigationEntry
	var index int64
	err = Run(ctx, Navigate(testdataDir+"/image.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	err = Run(ctx, NavigationEntries(&index, &entries))
	if err != nil {
		t.Fatal(err)
	}

	err = Run(ctx, Navigate(testdataDir+"/form.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	err = Run(ctx, NavigateToHistoryEntry(entries[index].ID))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var title string
	err = Run(ctx, Title(&title))
	if err != nil {
		t.Fatal(err)
	}
	if title != entries[index].Title {
		t.Errorf("expected title to be %s, instead title is: %s", entries[index].Title, title)
	}
}

func TestNavigateBack(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	err = Run(ctx, Navigate(testdataDir+"/form.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var exptitle string
	err = Run(ctx, Title(&exptitle))
	if err != nil {
		t.Fatal(err)
	}

	err = Run(ctx, Navigate(testdataDir+"/image.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	err = Run(ctx, NavigateBack())
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var title string
	err = Run(ctx, Title(&title))
	if err != nil {
		t.Fatal(err)
	}
	if title != exptitle {
		t.Errorf("expected title to be %s, instead title is: %s", exptitle, title)
	}
}

func TestNavigateForward(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	err = Run(ctx, Navigate(testdataDir+"/form.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	err = Run(ctx, Navigate(testdataDir+"/image.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var exptitle string
	err = Run(ctx, Title(&exptitle))
	if err != nil {
		t.Fatal(err)
	}

	err = Run(ctx, NavigateBack())
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	err = Run(ctx, NavigateForward())
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var title string
	err = Run(ctx, Title(&title))
	if err != nil {
		t.Fatal(err)
	}
	if title != exptitle {
		t.Errorf("expected title to be %s, instead title is: %s", exptitle, title)
	}
}

func TestStop(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	err = Run(ctx, Navigate(testdataDir+"/form.html"))
	if err != nil {
		t.Fatal(err)
	}

	err = Run(ctx, Stop())
	if err != nil {
		t.Fatal(err)
	}
}

func TestReload(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	err = Run(ctx, Navigate(testdataDir+"/form.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var exptitle string
	err = Run(ctx, Title(&exptitle))
	if err != nil {
		t.Fatal(err)
	}

	err = Run(ctx, Reload())
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var title string
	err = Run(ctx, Title(&title))
	if err != nil {
		t.Fatal(err)
	}
	if title != exptitle {
		t.Errorf("expected title to be %s, instead title is: %s", exptitle, title)
	}
}

func TestCaptureScreenshot(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	err = Run(ctx, Navigate(testdataDir+"/image.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var buf []byte
	err = Run(ctx, CaptureScreenshot(&buf))
	if err != nil {
		t.Fatal(err)
	}

	if len(buf) == 0 {
		t.Fatal("failed to capture screenshot")
	}
	//TODO: test image
}

/*func TestAddOnLoadScript(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	var scriptID page.ScriptIdentifier
	err = Run(ctx, AddOnLoadScript(`window.alert("TEST")`, &scriptID))
	if err != nil {
		t.Fatal(err)
	}

	err = Run(ctx, Navigate(testdataDir+"/form.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	if scriptID == "" {
		t.Fatal("got empty script ID")
	}
	// TODO: Handle javascript dialog.
}

func TestRemoveOnLoadScript(t *testing.T) {
	t.Parallel()

	var err error

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	var scriptID page.ScriptIdentifier
	err = Run(ctx, AddOnLoadScript(`window.alert("TEST")`, &scriptID))
	if err != nil {
		t.Fatal(err)
	}

	if scriptID == "" {
		t.Fatal("got empty script ID")
	}

	err = Run(ctx, RemoveOnLoadScript(scriptID))
	if err != nil {
		t.Fatal(err)
	}

	err = Run(ctx, Navigate(testdataDir+"/form.html"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)
}*/

func TestLocation(t *testing.T) {
	t.Parallel()

	var err error
	expurl := testdataDir + "/form.html"

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	err = Run(ctx, Navigate(expurl))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var urlstr string
	err = Run(ctx, Location(&urlstr))
	if err != nil {
		t.Fatal(err)
	}

	if urlstr != expurl {
		t.Fatalf("expected to be on form.html, got: %s", urlstr)
	}
}

func TestTitle(t *testing.T) {
	t.Parallel()

	var err error
	expurl, exptitle := testdataDir+"/image.html", "this is title"

	ctx, cancel := testAllocate(t, "")
	defer cancel()

	err = Run(ctx, Navigate(expurl))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	var title string
	err = Run(ctx, Title(&title))
	if err != nil {
		t.Fatal(err)
	}

	if title != exptitle {
		t.Fatalf("expected title to be %s, got: %s", exptitle, title)
	}
}
