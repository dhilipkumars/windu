package main

import (
	"fmt"
	"os"
	"strings"
	"sort"
)

type filestat struct {
	name string
	size int64
}

//du is disk utilization usually recursive from the said location.
func du(filePath string) ([]filestat, error) {

	var result []filestat


	stat, err := os.Lstat(filePath)
	if err != nil {
		return result, err
	}

	isSymLink := stat.Mode() & os.ModeSymlink

	if isSymLink != 0 {
		//If this is is a symbolic link
		//fmt.Printf("Skipping a symbolic link %s\n", filePath)
		return result, nil
	}


	if !stat.IsDir() {
		//If this is a file
		//fmt.Printf("Is not a directory\n")
		result = append(result, filestat{name: stat.Name(), size: stat.Size()})
		return result, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		if strings.Contains(err.Error(), "Access is denied") {
			//If the error is access denied just complain and skip it
			fmt.Fprintf(os.Stderr, "open(%s) err=%v\n", filePath, err)
			return result, nil
		}
		return result, err
	}
	defer file.Close()

	//fmt.Printf("Is a directory\n")
	filesInDir, err := file.Readdir(-1)
	if err != nil {
		return result, err
	}

	//fmt.Printf("Sub directories are =%v\n", filesInDir)

	for _, subFile := range filesInDir {
		//fmt.Printf("Processing %s isDir:%v\n", subFile.Name(),subFile.IsDir() )
		if subFile.IsDir() {
			var dirStat filestat
			subResult, err := du(filePath + "\\" + subFile.Name())
			if err != nil {
				return result, err
			}
			dirStat.name = subFile.Name()
			for _, dir := range subResult {
				dirStat.size += dir.size
			}
			result = append(result, dirStat)

		} else {
			result = append(result, filestat{subFile.Name(), subFile.Size()})
		}
	}

	return result, nil

}

func printUsage(args []string) {
	fmt.Printf("Error: Should be %s <Dir Path>", args[0])
}

type BySize []filestat

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].size > a[j].size }

func printResult(results []filestat) {
	maxSize := 0

	//Sort the results according to max size
	sort.Sort(BySize(results))

	for _, result := range results {
		sizeLen := len(fmt.Sprintf("%d", result.size))
		if maxSize < sizeLen {
			maxSize = sizeLen
		}
	}

	for _, result := range results {
		fmt.Printf("%*d %s\n", maxSize, result.size, result.name)
	}
}

func main() {
	//Comamnd line arguments
	if len(os.Args) != 2 {
		printUsage(os.Args)
		os.Exit(1)
	}
	flist, err := du(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	printResult(flist)

}
