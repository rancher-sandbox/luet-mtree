/*
Copyright © 2021 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package action

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vbatts/go-mtree"
	"os"
)

type checkAction struct {
	target         string
	validationFile string
	format         string
}

func NewCheckAction(t string, v string, f string) *checkAction {
	return &checkAction{target: t, validationFile: v, format: f}
}

func (action checkAction) Run() error {
	// If its not a dir, try to uncompress
	info, _ := os.Stat(action.target)
	if !info.IsDir() {
		tmpDir, _ := os.MkdirTemp("", "luet-mtree")
		defer os.RemoveAll(tmpDir)
		newTarget, err := unTar(action.target, tmpDir)
		if err != nil { return err }
		action.target = newTarget
	}

	spec := &mtree.DirectoryHierarchy{}
	stateDh := &mtree.DirectoryHierarchy{}
	var err error
	var excludes []mtree.ExcludeFunc
	// excludeEmptyFiles is an ExcludeFunc for excluding all files with 0 size
	var excludeEmptyFiles = func(path string, info os.FileInfo) bool {
		if info.Size() == 0{
			return true
		}
		return false
	}
	excludes = append(excludes, excludeEmptyFiles)
	var res []mtree.InodeDelta

	fh, err := os.Open(action.validationFile)
	if err != nil {
		return err
	}
	spec, err = mtree.ParseSpec(fh)
	err = fh.Close()
	if err != nil {
		return err
	}

	specKeywords := spec.UsedKeywords()
	stateKeyworks := spec.UsedKeywords()

	stateDh, err = mtree.Walk(action.target, excludes, stateKeyworks, nil)
	res, err = mtree.Compare(spec, stateDh, specKeywords)
	if err != nil {
		return err
	}

	out := formats[action.format](res)
	if _, err := os.Stdout.Write([]byte(out)); err != nil {
		return err
	}

	for _, diff := range res {
		if diff.Type() == mtree.Modified {
			return errors.New("validation failed")
		}
	}
	return nil
}

var formats = map[string]func([]mtree.InodeDelta) string{
	// Outputs the errors in the BSD format.
	"bsd": func(d []mtree.InodeDelta) string {
		var buffer bytes.Buffer
		for _, delta := range d {
			_, _ = fmt.Fprintln(&buffer, delta)
		}
		return buffer.String()
	},

	// Outputs the full result struct in JSON.
	"json": func(d []mtree.InodeDelta) string {
		var buffer bytes.Buffer
		if err := json.NewEncoder(&buffer).Encode(d); err != nil {
			panic(err)
		}
		return buffer.String()
	},

	// Outputs only the paths which failed to validate.
	"path": func(d []mtree.InodeDelta) string {
		var buffer bytes.Buffer
		for _, delta := range d {
			if delta.Type() == mtree.Modified {
				_, _ = fmt.Fprintln(&buffer, delta.Path())
			}
		}
		return buffer.String()
	},
}
