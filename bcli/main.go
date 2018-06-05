package main

import (
	"fmt"
)

func main() {
	err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	cli := NewCLI()
	cli.Run()
}
