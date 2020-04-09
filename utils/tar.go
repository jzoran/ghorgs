// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.
//
// (inspired by snippet from https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726)

package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func TarGz(dest, src string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("Could not create %s. Error! %s", dest, err.Error())
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.Walk(src, func(f string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Could not traverse '%s'. Error! %s",
				f, err.Error())
		}

		hf := f
		if fi.Mode()&os.ModeSymlink != 0 {
			hf, err = os.Readlink(f)
			if err != nil {
				return fmt.Errorf("Could not read symlink '%s'. Error! %s",
					f, err.Error())
			}
		}

		h, err := tar.FileInfoHeader(fi, hf)
		if err != nil {
			return fmt.Errorf("Could not generate tar header for '%s'. Error! %s",
				dest, err.Error())
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		h.Name = filepath.ToSlash(f)
		err = tw.WriteHeader(h)
		if err != nil {
			return fmt.Errorf("Could not write header for '%s'. Error! %s",
				dest, err.Error())
		}
		// if a regular file, write file content
		if fi.Mode().IsRegular() {
			data, err := os.Open(f)
			if err != nil {
				return fmt.Errorf("Could not open '%s'. Error! %s", f, err.Error())
			}
			_, err = io.Copy(tw, data)
			if err != nil {
				return fmt.Errorf("Could not write '%s' to '%s'. Error! %s",
					f, dest, err.Error())
			}
		}
		return nil
	})
}
