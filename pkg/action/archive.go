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

	tarReader := tar.NewReader(r)

	for {
		header, err := tarReader.Next()
		// Remember to access any values from header below this error check, otherwise you could be accessing an empty
		// header and provoke a runtime error!
		switch {
			// if no more files are found return as we have our extracted dir
			case err == io.EOF:
				return tmpDir, nil
			// if error or header nil, not a tar archive
			case err != nil || header == nil:
				fmt.Printf("Error reading file: %v (Not a tar archive?)\n", err)
				return "", err
		}

		fileInfo := header.FileInfo()
		targetDir = filepath.Join(tmpDir, header.Name)

		if fileInfo.IsDir() {
			if _, err := os.Stat(targetDir); err != nil {
				_ = os.MkdirAll(targetDir, fileInfo.Mode().Perm())
			}
		} else {
			f, err := os.OpenFile(targetDir, os.O_CREATE|os.O_RDWR, fileInfo.Mode().Perm())
			if err != nil {
				return "", fmt.Errorf("failure creating file %s: %s\n", header.Name, err)
			}
			n, err := io.Copy(f, tarReader)
			if err != nil {
				return "", fmt.Errorf("failure to copy to file %s: %s\n", header.Name, err)
			}
			err = f.Close()
			if err != nil{
				return "", fmt.Errorf("failure closing file %s: %s\n", header.Name, err)
			}

			if n != fileInfo.Size() {
				return "", fmt.Errorf("size extracted for %s differs from size wanted: %d -> %d\n", header.Name, n, fileInfo.Size())
			}
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