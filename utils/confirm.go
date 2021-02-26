//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package utils

import "fmt"

// GetUserConfirmation displays prompt and acquires user input
// required for some commands.
func GetUserConfirmation() bool {
	if Debug.DryRun {
		fmt.Println("This is a dry-run. No data will actually be modified on GitHub and only the " +
			"commands will be print out.")
	} else {
		fmt.Println("ATTN: This action may be irreversible.")
	}
	fmt.Println("Are you sure you want to continue? (y/N):")
	var yes string
	n, err := fmt.Scanf("%s", &yes)
	if err != nil {
		fmt.Println(n)
		fmt.Println(err.Error())
		return false
	}

	return yes == "d" ||
		yes == "da" ||
		yes == "D" ||
		yes == "Da" ||
		yes == "j" ||
		yes == "ja" ||
		yes == "J" ||
		yes == "Ja" ||
		yes == "y" ||
		yes == "yes" ||
		yes == "yep" ||
		yes == "Y" ||
		yes == "Yes" ||
		yes == "Sure thing, mate! Please do carry on."
}
