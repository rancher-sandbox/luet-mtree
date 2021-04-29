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
	"os"
)

type generateAction struct {
	target string
	outputFile string
	keywords []string
}

func NewGenerateAction(t string, o string, k []string) *generateAction {
	return &generateAction{target: t, outputFile: o, keywords: k}
}


func (action generateAction) Run() error {
	stateDh := &mtree.DirectoryHierarchy{}
	var excludes []mtree.ExcludeFunc
	var err error

	fh := os.Stdout
	if action.outputFile != "" {
		fh, err = os.Create(action.outputFile)
		if err != nil {
			return err
		}
	}

	// TODO(itxaka): Do we want all keywords or its enough with just the sha?
	currentKeywords := mtree.DefaultKeywords[:]
	if len(action.keywords) > 0 {
		currentKeywords = mtree.ToKeywords(action.keywords)
	}
	fmt.Print(currentKeywords)

	currentKeywords = append(currentKeywords, "sha256digest")
	stateDh, err = mtree.Walk(action.target, excludes, currentKeywords, nil)
	if err != nil {
		return err
	}

	_, err = stateDh.WriteTo(fh)
	if err != nil {
		return err
	}

	return nil
}