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
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"io"
	"os"
	"path/filepath"
)

const (
	GZIP = "application/gzip"
	ZSTD = "application/zstd"
	NONE = "none"
)


func unTar(target string, tmpDir string) (string, error) {
	targetDir := ""
	original, _ := os.Open(target)
	buf := bytes.Buffer{}

	if compressType(target) == GZIP {
		r, _ := gzip.NewReader(original)
		_, _ = r.WriteTo(&buf)
		fmt.Printf("Found gzip compression\n")
	} else if compressType(target) == ZSTD{
		r, _ := zstd.NewReader(original)
		_, _ = r.WriteTo(&buf)
		fmt.Printf("Found zstd compression\n")
	} else {
		_, _ = io.Copy(&buf, original)
		fmt.Printf("Found no compression\n")
	}

	tarReader := tar.NewReader(&buf)

	for {
		header, err := tarReader.Next()

		switch {
			// if no more files are found return as we have our extracted dir
			case err == io.EOF:
				return tmpDir, nil
			// if error or header nil, not a tar archive
			case err != nil || header == nil:
				fmt.Printf("Error reading file: %v (Not a tar archive?)\n", err)
				return "", err
		}

		targetDir = filepath.Join(tmpDir, header.Name)
		switch header.Typeflag {
			// Dir
			case tar.TypeDir:
				if _, err := os.Stat(targetDir); err != nil {
					_ = os.MkdirAll(targetDir, os.FileMode(header.Mode))
				}

			// file
			case tar.TypeReg:
				f, _ := os.OpenFile(targetDir, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				_, _ = io.Copy(f, tarReader)
				_ = f.Close()
			}
		}
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