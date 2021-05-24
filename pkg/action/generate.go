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
	"fmt"
	"github.com/vbatts/go-mtree"
	"io"
	"io/ioutil"
	"os"
)

type generateAction struct {
	target     string
	outputFile string
	keywords   []string
}

func NewGenerateAction(t string, o string, k []string) *generateAction {
	return &generateAction{target: t, outputFile: o, keywords: k}
}

func (action generateAction) Run() error {
	// If its not a dir, try to uncompress
	info, _ := os.Stat(action.target)
	stateDh := &mtree.DirectoryHierarchy{}
	var excludes []mtree.ExcludeFunc

	var err error

	// excludeEmptyFiles is an ExcludeFunc for excluding all files with 0 size
	//var excludeEmptyFiles = func(path string, info os.FileInfo) bool {
	//	if info.Size() == 0{
	//return true
	//}
	//	return false
	//}
	//excludes = append(excludes, excludeEmptyFiles)

	var infoFiles = func(path string, info os.FileInfo) bool {
		fmt.Printf("%s -> %s with size %s\n", path, info.Name(), info.Size())
		return false
	}
	excludes = append(excludes, infoFiles)
	fh := os.Stdout
	if action.outputFile != "" {
		fh, err = os.OpenFile(action.outputFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
	}

	// Ignore time because luet tars files with the docker lib and that truncates the time to seconds only
	// TODO(itxaka) we may be able to use tar_time ?
	currentKeywords := []mtree.Keyword{
		"type",
		"xattrs",
		"sha512digest",
	}

	if len(action.keywords) > 0 {
		for _, k := range action.keywords {
			if !mtree.InKeywordSlice(mtree.Keyword(k), currentKeywords) {
				currentKeywords = append(currentKeywords, mtree.Keyword(k))
			}
		}
	}

	if !info.IsDir() {
		uncompressedTar, err := unCompress(action.target)
		ts := mtree.NewTarStreamer(uncompressedTar, excludes, currentKeywords)
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
		stateDh, err = mtree.Walk(action.target, excludes, currentKeywords, nil)
		if err != nil {
			return err
		}
	}

	_, err = stateDh.WriteTo(fh)
	if err != nil {
		return err
	}

	return nil
}
