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
	"github.com/gabriel-vasile/mimetype"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"io"
	"os"
)

const (
	GZIP = "application/gzip"
	ZSTD = "application/zstd"
	NONE = "none"
)


func unCompress(target string) (io.Reader, error) {
	original, _ := os.Open(target)
	var r io.Reader

	if compressType(target) == GZIP {
		r, _ = gzip.NewReader(original)
		fmt.Printf("Found gzip compression\n")
	} else if compressType(target) == ZSTD{
		r, _ = zstd.NewReader(original)
		fmt.Printf("Found zstd compression\n")
	} else {
		r = original
		fmt.Printf("Found no compression\n")
	}

	return r, nil
}


func compressType(file string) string {
	mime, _ := mimetype.DetectFile(file)
	switch mime.String() {
	case GZIP:
		return GZIP
	case ZSTD:
		return ZSTD
	default:
		return NONE
	}
}