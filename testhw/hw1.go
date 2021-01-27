package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("work")
	out := os.Stdout
	fmt.Printf("%#v", out)
	filepa := (filepath.Dir(""))
	fmt.Printf("%s", filepa)
	files, err := FilePathWalkDir(".")
	if err != nil {
		panic(err)
	}
	// fmt.Println(files)
	var a []string
	for _, value := range files {
		a = append(a, value)
	}
	fmt.Print(a)
	hel := "hello world fuck"
	helSplit := strings.Fields(hel)
	fmt.Println(hel)
	fmt.Println(helSplit)
	fmt.Println(helSplit[len(helSplit)-1])
	for _, j := range helSplit {
		fmt.Println(j + "/")
	}
	treeMap, parents := createMaps(a)
	fmt.Println(treeMap)
	fmt.Println(parents)
	var tree string
	var listBool map[int]bool = map[int]bool{
		0: false,
	}
	for n, parent := range parents {
		if n != len(parents)-1 {
			tree += "├───" + parent + "\n"
			listBool[0] = false
			tree = paintTree(treeMap, parent, tree, 1, listBool)
		} else {
			tree += "└───" + parent + "\n"
			listBool[0] = true
			tree = paintTree(treeMap, parent, tree, 1, listBool)
		}

	}
	fmt.Print(tree)

}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			files = append(files, path)
		}
		// fmt.Printf("%s\n", info.Name())
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
			parents = append(parents, pathDirList[0])
		case countDirs > 1:
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
					tree += "|\t"
				}
			}
			tree += symbol + child + "\n"
			tree = paintTree(treeMap, child, tree, label+1, lastBool)
		}
	}
	return tree
}
