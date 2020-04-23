// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package utils

type DebugConfiguration struct {
	Verbose bool
	DryRun  bool
}

var Debug = DebugConfiguration{false, false}
