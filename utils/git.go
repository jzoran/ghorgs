//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
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
		return fmt.Errorf("`git clone` failed with %s\n", err.Error())
	}

	return nil
}

func Url(rawurl, user, pass string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("Project url (%s) error. %s", rawurl, err.Error())
	}

	if user != "" && pass != "" {
		u.User = url.UserPassword(user, pass)
	}
	return u.String(), nil
}
