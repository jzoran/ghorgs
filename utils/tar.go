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

func TargzVerify(dest, src string) error {
	ar, err := os.Open(dest)
	if err != nil {
		return fmt.Errorf("Could not open '%s' for reading. Error! %s",
			dest, err.Error())
	}
	defer ar.Close()

	fimap := make(map[string]os.FileInfo)
	err = filepath.Walk(src, func(f string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Could not traverse '%s'. Error! %s",
				f, err.Error())
		}
		fimap[f] = fi
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error while traversing '%s'. Error! %s", src, err.Error())
	}

	gr, err := gzip.NewReader(ar)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return fmt.Errorf("Could not uncompress archive '%s'. Error! %s",
			dest, err.Error())
	}
	tr := tar.NewReader(gr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Header check error for archive '%s'. Error! %s",
				dest, err.Error())
		}

		fi, ok := fimap[h.Name]
		if ok {
			if fi.Mode().IsRegular() && fi.Size() != h.Size {
				return fmt.Errorf("Incorrect size of '%s' in archive '%s'. Expected '%d',"+
					" but found '%d'", h.Name, dest, fi.Size(), h.Size)
			}
			if fi.Mode() != h.FileInfo().Mode() {
				return fmt.Errorf("Incorrect stat entry for '%s' in archive '%s'. ",
					h.Name, dest)
			}
			if h.FileInfo().Mode()&os.ModeSymlink != 0 {
				link, err := os.Readlink(h.Name)
				if err != nil {
					return fmt.Errorf("Could not read symlink '%s' in archive '%s'. Error! %s",
						h.Name, dest, err.Error())

				}
				if link != h.Linkname {
					return fmt.Errorf("Incorrect symlink of '%s' in archive '%s'. Expected '%s',"+
						" but found' %s'.",
						h.Name, dest, link, h.Linkname)
				}
			}
			delete(fimap, h.Name)
		}
	}
	if len(fimap) > 0 {
		return fmt.Errorf("Incorrect archive '%s'. Missing files from '%s': %v.",
			dest, src, keysOf(fimap))
	}

	return nil
}

func keysOf(m map[string]os.FileInfo) []string {
	keys := make([]string, 0)
	for key, _ := range m {
		keys = append(keys, key)
	}
	return keys
}
