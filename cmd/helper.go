package cmd

import "fmt"

func console(err error, value interface{}) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(" %v\n", value)
	}
}
