package main

import (
	"fmt"
	"os"
)

type filestat struct {
	name string
	size int64
}

//du is disk utilization usually recursive from the said location.
func du(filePath string) ([]filestat, error) {

	var result []filestat

	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return result, err
	}

	if !stat.IsDir() {
		//If this is a file
		//fmt.Printf("Is not a directory\n")
		result = append(result, filestat{name: stat.Name(), size: stat.Size()})
		return result, nil
	}

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

func printResult(results []filestat) {
	maxSize := 0

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
