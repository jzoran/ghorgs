// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

func GitClone(url string, out string, name string) error {
	// assumes url and dest are valid
	dest := path.Join(out, name)
	cmd := exec.Command("git", "clone", url, dest)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("`git clone` failed with %s\n", err)
	}

	return nil
}
