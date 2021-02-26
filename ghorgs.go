//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package main

import (
	"fmt"
	"ghorgs/cmd"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
