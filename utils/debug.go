//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package utils

type DebugConfiguration struct {
	Verbose bool
	DryRun  bool
}

var Debug = DebugConfiguration{false, false}
