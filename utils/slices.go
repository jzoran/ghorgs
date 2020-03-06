// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package utils

func StringInSlice(s string, list []string) bool {
	for _, item := range list {
		if s == item {
			return true
		}
	}
	return false
}
