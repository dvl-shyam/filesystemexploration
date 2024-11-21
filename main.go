package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: fileman <command> [options]")
		return
	}
	fmt.Println(os.Args)
	switch os.Args[1] {
	case "list":
		ListFiles()
	case "search":
		searchFile()
	case "copy":
	case "move":
	case "delete":
	default:
		fmt.Println("Unknown Command", os.Args[1])
	}

}
