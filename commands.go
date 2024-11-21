package main

import (
	"flag"
	"fmt"
	"os"
)

func ListFiles() {

	flagSet := flag.NewFlagSet("list", flag.ExitOnError)
	path := flagSet.String("path", ".", "path to list files / directory")
	flagSet.Parse(os.Args[2:])

	files, err := os.ReadDir(*path)
	if err != nil {
		fmt.Println("Error Reading Directory", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println("[DIR]", file.Name())
		} else {
			fmt.Println("[FILE]", file.Name())
		}
	}

}

func searchFile(){
	flagSet := flag.NewFlagSet("search", flag.ExitOnError)
	name := flagSet.String("name", ".", "search query name")
	path := flagSet.String("path", ".", "Path to Diretory")
	flagSet.Parse(os.Args[2:])

	files, err := os.ReadDir(*path)
	if err != nil{
		fmt.Println("error readidng directory", err)
	}

	if *name == "" {
		fmt.Println("Error: Please provide a name or extension to search for.")
		return
	}

	for _, file := range files {
		if file.Name() == *name {
			fmt.Println("File Exist")
		} else {
			fmt.Println("File Doesnt exist")
		}
	}

}
