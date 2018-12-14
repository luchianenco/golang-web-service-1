package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, showFiles bool) error {
	isPathDir, err := checkPathDir(path)
	if err != nil {
		panic(err)
	}
	if !isPathDir {
		panic("The given path is not a directory")
	}

	printEntries(out, path, showFiles, "")

	return err
}

func printEntries(out io.Writer, path string, showFiles bool, ident string) error {
	entries := getEntries(path, showFiles)
	isLast := false
	prefix := getPrefix(isLast)
	for key, entry := range entries {
		if key == len(entries)-1 {
			isLast = true
			prefix = getPrefix(isLast)
		}

		printName(out, entry, ident+prefix)

		if entry.IsDir() {
			subPath := filepath.Join(path, entry.Name())
			subIdent := getIdent(ident, isLast)
			printEntries(out, subPath, showFiles, subIdent)
		}
	}

	return nil
}

func getEntries(path string, showFiles bool) []os.FileInfo {
	var filteredEntries []os.FileInfo
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !showFiles && !entry.IsDir() {
			continue
		}
		filteredEntries = append(filteredEntries, entry)

	}

	return filteredEntries
}

func printName(out io.Writer, entry os.FileInfo, prefix string) {
	size := getEntrySize(entry)
	name := prefix + entry.Name() + size

	fmt.Fprintln(out, name)
}

func getPrefix(currentLast bool) string {
	prefix := "├───"

	if currentLast {
		prefix = "└───"
	}

	return prefix
}

func getEntrySize(entry os.FileInfo) string {
	if entry.IsDir() {
		return ""
	}

	if entry.Size() == 0 {
		return " (empty)"
	}

	return fmt.Sprintf(" (%db)", entry.Size())
}

func getIdent(ident string, isLast bool) string {
	identPrefix := ""

	if isLast {
		identPrefix = ident + "	"
	} else {
		identPrefix = ident + "│	"
	}

	return identPrefix
}

func checkPathDir(path string) (bool, error) {
	pathInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return pathInfo.IsDir(), err
}
