package action

import (
	"bytes"
	"github.com/vbatts/go-mtree"
	"os"
	"testing"
)

// captureOutputAndParse hijacks stdout and uses that output to try to return an
// mtree.DirectoryHierarchy
// This allows us to parse the mtree output from generate so we can easily parse keywords used for example
func captureOutputAndParse(t *testing.T, f func()) CustomDirectoryHierarchy {
	// Overwrite stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	origStdout := os.Stdout
	os.Stdout = w
	// Call function
	f()
	// Store the output
	buf := make([]byte, 2048)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	// Restore stdout
	os.Stdout = origStdout

	// Parse output
	res := bytes.NewReader(buf[:n])
	directoryHierarchy, err := mtree.ParseSpec(res)
	if err != nil {
		t.Fatal(err)
	}
	cdh := CustomDirectoryHierarchy{directoryHierarchy}

	return cdh
}

type CustomDirectoryHierarchy struct {
	*mtree.DirectoryHierarchy
}

// getRealKeywords gets the "real" keywords used to generate the file
// The current method from mtree checks for the /set values and those add extra default keywords in them
// for example uid and gid
// They are not used to check the mtree values, as those are not stored on the mtree output but its kind of
// confusing that those are not the real keywords used for the actual files so we add this simple method
// So we can check them safely in testing
func (cdh CustomDirectoryHierarchy) getRealKeywords() []mtree.Keyword {
	var usedkeywords []mtree.Keyword
	for _, e := range cdh.Entries {
		switch e.Type {
		case mtree.FullType, mtree.RelativeType:
			if e.Type != mtree.SpecialType {
				kvs := e.Keywords
				for _, kv := range kvs {
					kw := kv.Keyword()
					if !mtree.InKeywordSlice(kw, usedkeywords) {
						usedkeywords = append(usedkeywords, mtree.KeywordSynonym(string(kw)))
					}
				}
			}
		}
	}
	return usedkeywords
}
