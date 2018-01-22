package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type ByName []os.FileInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func dirTree(output io.Writer, path string, printFiles bool) error {
	return readPath(path, "", printFiles, output)
}

func readPath(path string, prefix string, printFiles bool, output io.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	files, err := file.Readdir(0)
	sort.Sort(ByName(files))

	if printFiles == false {
		files = getDirs(files)
	}

	len := len(files)
	for index, element := range files {
		isDir := element.IsDir()
		isLast := len == index+1

		fmt.Fprint(output, getCurrentPrefix(prefix, isLast)+element.Name()+getSize(element, isDir), "\n")

		if isDir {
			readPath(path+"/"+element.Name(), prefix+getDirPrefix(isLast), printFiles, output)
		}
	}

	return nil
}

func getDirs(files []os.FileInfo) []os.FileInfo {
	dirs := make([]os.FileInfo, 0)
	for _, element := range files {
		if element.IsDir() {
			dirs = append(dirs, element)
		}
	}
	return dirs
}

func getCurrentPrefix(prefix string, isLast bool) string {

	filePrefix := "├───"
	if isLast {
		filePrefix = "└───"
	}
	return prefix + filePrefix
}

func getDirPrefix(isLast bool) string {
	if isLast {
		return "	"
	} else {
		return "│	"
	}
}

func getSize(fileInfo os.FileInfo, isDir bool) string {
	size := ""
	if isDir == false {
		sizeB := int(fileInfo.Size())
		if sizeB == 0 {
			size = "empty"
		} else {
			size = strconv.Itoa(sizeB) + "b"
		}
		size = " (" + size + ")"
	}
	return size
}

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
