package action

import (
	"github.com/stretchr/testify/assert"
	"github.com/vbatts/go-mtree"
	"testing"
)



func TestGenerateExpandKeywords(t *testing.T) {
	output := captureOutputAndParse(t, func() {
		action := NewGenerateAction("testdata/checkfiles/", "", []string{"size"}, make([]string, 0))
		err := action.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("size"))
	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("sha512digest"))  // We add this as default
	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("type"))  // We add this as default
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("uid"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("gid"))

	output = captureOutputAndParse(t, func() {
		action := NewGenerateAction("testdata/checkfiles/", "", []string{"size", "mode"}, make([]string, 0))
		err := action.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("size"))
	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("mode"))
	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("sha512digest"))  // We add this as default
	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("type"))  // We add this as default
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("uid"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("gid"))
}

func TestGenerateOverrideKeywords(t *testing.T) {
	output := captureOutputAndParse(t, func() {
		action := NewGenerateAction("testdata/checkfiles/", "", []string{}, []string{"size"})
		err := action.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("size"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("sha512digest"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("type"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("uid"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("gid"))

	output = captureOutputAndParse(t, func() {
		action := NewGenerateAction("testdata/checkfiles/", "", []string{}, []string{"size", "mode"})
		err := action.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("size"))
	assert.Contains(t, output.getRealKeywords(), mtree.Keyword("mode"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("sha512digest"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("type"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("uid"))
	assert.NotContains(t, output.getRealKeywords(), mtree.Keyword("gid"))
}

