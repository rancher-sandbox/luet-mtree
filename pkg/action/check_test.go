package action

import (
	"testing"
)

func TestCheck(t *testing.T) {
	action := NewCheckAction("testdata/checkfiles/", "testdata/checkfiles.sum", "bsd", make([]string, 0))
	out, err := action.Run()
	if out != "" {
		t.Fatalf("Expected empty output, got: %v", out)
	}
	if err != nil {
		t.Fatal(err)
	}
}
