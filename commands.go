package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func searchFile() {
	flagSet := flag.NewFlagSet("search", flag.ExitOnError)
	name := flagSet.String("name", ".", "search query name")
	path := flagSet.String("path", ".", "Path to Diretory")
	flagSet.Parse(os.Args[2:])

	if *name == "" {
		fmt.Println("Error: Please provide a name or extension to search for.")
		return
	}

	found := false

	err := filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), *name) {
			fmt.Println(path)
			found = true
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error searching files:", err)
	} else if !found {
		fmt.Println("No file found")
	}

}

func copyDir(srcDir, destDir string) error {
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	err = filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if srcPath == srcDir {
			return nil
		}
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		return copyFile(srcPath, destPath)
	})

	if err != nil {
		return fmt.Errorf("error while copying directory contents: %w", err)
	}

	return nil
}
func copyFile(src, dest string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("error while copying file content: %w", err)
	}

	return nil
}

func copying() {
	flagSet := flag.NewFlagSet("copy", flag.ExitOnError)
	src := flagSet.String("src", ".", "source")
	dest := flagSet.String("dest", ".", "Destination")
	flagSet.Parse(os.Args[2:])

	srcStat, err := os.Stat(*src)
	if err != nil {
		fmt.Println("Error occured in file 71", err)
		return
	}
	destStat, err := os.Stat(*dest)
	if err != nil {
		fmt.Println("Error occurred in destination stat:", err)
		return
	}
	if !srcStat.IsDir() {
		if destStat.IsDir() {
			dest = flag.String("dest", filepath.Join(*dest, filepath.Base(*src)), "Destination file path")
		}
		err := copyFile(*src, *dest)
		if err != nil {
			fmt.Println("Error occurred while copying file:", err)
		} else {
			fmt.Println("File copied successfully.")
		}
		return
	}

	if srcStat.IsDir() {
		if !destStat.IsDir() {
			fmt.Println("Error: destination must be a directory when copying a directory.")
			return
		}

		err := copyDir(*src, *dest)
		if err != nil {
			fmt.Println("Error occurred while copying directory:", err)
		} else {
			fmt.Println("Directory copied successfully!")
		}
	}

	fmt.Println("File copied successfully.")
}

func moveFile(src, dest string) error {
	err := copyFile(src, dest)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("failed to remove source file after moving: %w", err)
	}

	return nil
}

func moveDir(srcDir, destDir string) error {
	srcBase := filepath.Base(srcDir)

	newDestDir := filepath.Join(destDir, srcBase)

	err := os.MkdirAll(newDestDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	err = filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}
		destPath := filepath.Join(newDestDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		return copyFile(srcPath, destPath)
	})

	if err != nil {
		return fmt.Errorf("error while copying directory contents: %w", err)
	}
	err = os.RemoveAll(srcDir)
	if err != nil {
		return fmt.Errorf("failed to remove source directory after moving: %w", err)
	}

	return nil
}

func move() {
	flagSet := flag.NewFlagSet("move", flag.ExitOnError)
	src := flagSet.String("src", ".", "source")
	dest := flagSet.String("dest", ".", "Destination")
	flagSet.Parse(os.Args[2:])

	srcStat, err := os.Stat(*src)
	if err != nil {
		fmt.Println("Error occurred in source stat:", err)
		return
	}

	if !srcStat.IsDir() {
		destPath := *dest
		if stat, err := os.Stat(*dest); err == nil && stat.IsDir() {
			destPath = filepath.Join(*dest, filepath.Base(*src))
		}

		err := moveFile(*src, destPath)
		if err != nil {
			fmt.Println("Error occurred while moving file:", err)
		} else {
			fmt.Println("File moved successfully.")
		}
		return
	}

	err = moveDir(*src, *dest)
	if err != nil {
		fmt.Println("Error occurred while moving directory:", err)
	} else {
		fmt.Println("Directory moved successfully!")
	}
}

func delete() {
	flagSet := flag.NewFlagSet("delete", flag.ExitOnError)
	path := flagSet.String("path", ".", "Path to file or directory to delete")
	flagSet.Parse(os.Args[2:])

	stat, err := os.Stat(*path)
	if err != nil {
		fmt.Println("Error occurred while accessing the path:", err)
		return
	}

	if stat.IsDir() {
		err := os.RemoveAll(*path)
		if err != nil {
			fmt.Printf("failed to delete directory: %s", err)
			return
		} else {
			fmt.Println("Directory deleted successfully!")
		}
	} else {
		err := os.Remove(*path)
		if err != nil {
			fmt.Printf("failed to delete file: %s", err)
			return
		} else {
			fmt.Println("File deleted successfully!")
		}
	}
}
