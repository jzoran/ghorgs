// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package utils

import "fmt"

func GetUserConfirmation() bool {
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
