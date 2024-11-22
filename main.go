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
	switch os.Args[1] {
	case "list":
		ListFiles()
	case "search":
		searchFile()
	case "copy":
		copying()
	case "move":
		move()
	case "delete":
		delete()
	default:
		fmt.Println("Unknown Command", os.Args[1])
	}

}
