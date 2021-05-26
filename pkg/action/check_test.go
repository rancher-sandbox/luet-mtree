package action

import "testing"

func TestCheck(t *testing.T) {
	action := NewCheckAction("testdata/checkfiles/", "testdata/checkfiles.sum", "bsd", make([]string, 0))
	err := action.Run()
	if err != nil {
		t.Fatal(err)
	}
}
