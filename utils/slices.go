//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package utils

// StringInSlice looks for a string `needle` in slice `haystack` and returns true
// if successful and false otherwise.
func StringInSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if needle == item {
			return true
		}
	}
	return false
}
