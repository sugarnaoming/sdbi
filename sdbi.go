package main

import (
	"fmt"
	"sdbi/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
