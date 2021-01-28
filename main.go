package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func dirTree(output io.Writer, path string, printFiles bool) error {
	files, err := FilePathWalkDir(path, printFiles)
	fmt.Println(files)
	if err != nil {
		panic(err)
	}
	var allPathElements []string
	for _, value := range files {
		allPathElements = append(allPathElements, value)
	}
	treeMap, parents := createMaps(allPathElements)
	var tree string
	var listBool map[int]bool = map[int]bool{
		0: false,
	}
	for n, parent := range parents {
		if n != len(parents)-1 {
			tree = paintTree(treeMap, parent, tree, 0, listBool)
		} else {
			tree = paintTree(treeMap, parent, tree, 0, listBool)
			tree = strings.TrimSuffix(tree, "\n")

		}

	}
	fmt.Fprintln(output, tree)
	return nil
}

func FilePathWalkDir(root string, printFiles bool) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if printFiles {
			size := info.Size()
			if size == 0 {
				path += " (empty)"
			} else if !info.IsDir() {
				strSize := strconv.FormatInt(size, 10)
				path += " (" + strSize + "b)"
			}
			files = append(files, path)
		} else if info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func createMaps(paths []string) (map[string][]string, []string) {
	workMaps := make(map[string][]string)
	var parents []string
	for _, path := range paths {
		pathDirList := strings.Split(path, "/")
		countDirs := len(pathDirList)
		switch {
		case countDirs == 1:
			if pathDirList[0] == "." || pathDirList[0] == ".git" || pathDirList[0] == ".DS_Store (6148b)" {
				continue
			}
			parents = append(parents, pathDirList[0])
		case countDirs > 1:
			if pathDirList[len(pathDirList)-1] == ".DS_Store (6148b)" {
				continue
			}
			workMaps[pathDirList[len(pathDirList)-2]] = append(workMaps[pathDirList[len(pathDirList)-2]], pathDirList[len(pathDirList)-1])

		}
	}
	return workMaps, parents

}

func paintTree(treeMap map[string][]string, parent string, tree string, label int, lastBool map[int]bool) string {
	childs, isParent := treeMap[parent]
	if isParent {
		countChild := len(childs)
		for n, child := range childs {
			symbol := ""
			if countChild > 1 && n != countChild-1 {
				symbol = "├───"
				lastBool[label] = false
			} else if countChild == 1 || n == countChild-1 {
				symbol = "└───"
				lastBool[label] = true
			}
			for i := 0; i < label; i++ {
				if lastBool[i] {
					tree += "\t"
				} else {
					tree += "│\t"
				}
			}
			tree += symbol + child + "\n"
			tree = paintTree(treeMap, child, tree, label+1, lastBool)
		}
	}
	return tree
}
