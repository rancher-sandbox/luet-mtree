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
	"github.com/itxaka/luet-mtree/pkg/log"
	"github.com/vbatts/go-mtree"
	"io"
	"io/ioutil"
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
	log.Log(fmt.Sprintf("Checking %s against validation file %s", action.target, action.validationFile))
	spec := &mtree.DirectoryHierarchy{}
	stateDh := &mtree.DirectoryHierarchy{}
	var err error
	var excludes []mtree.ExcludeFunc
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

	stateKeyworks := spec.UsedKeywords()

	// If its not a dir, try to uncompress
	info, _ := os.Stat(action.target)
	if !info.IsDir() {
		uncompressedTar, err := unCompress(action.target)
		ts := mtree.NewTarStreamer(uncompressedTar, excludes, stateKeyworks)
		if _, err := io.Copy(ioutil.Discard, ts); err != nil && err != io.EOF {
			return err
		}
		if err := ts.Close(); err != nil {
			return err
		}
		stateDh, err = ts.Hierarchy()
		if err != nil {
			return err
		}

	} else {
		stateDh, err = mtree.Walk(action.target, excludes, stateKeyworks, nil)
		if err != nil { return err }
	}

	res, err = mtree.Compare(spec, stateDh, stateKeyworks)
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
