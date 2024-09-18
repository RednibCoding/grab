package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var skippedFiles int
var skippedDirs = make(map[string]int)

const version = "1.0.1"

// isBinary checks the first 1024 bytes of a file to see if it's a binary file
func isBinary(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	// Heuristic: if there's a null byte in the first 1024 bytes, treat as binary
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true
		}
	}

	return false
}

// searchInFile reads a file line by line and searches for the searchString
func searchInFile(filePath string, searchString string, caseSensitive bool, results chan<- string, skippedDirs map[string]int) {
	defer wg.Done()

	// Skip binary files
	if isBinary(filePath) {
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		// Track skipped files and directories
		skippedFiles++
		dir := filepath.Dir(filePath)
		skippedDirs[dir]++
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineNumber := 1
	for {
		line, err := reader.ReadString('\n') // Read until the next newline character
		if err != nil && err != io.EOF {
			break
		}

		// Perform case-insensitive search if needed
		if !caseSensitive {
			line = strings.ToLower(line)
			searchString = strings.ToLower(searchString)
		}

		if strings.Contains(line, searchString) {
			column := strings.Index(line, searchString) + 1
			result := fmt.Sprintf("%s:%d:%d", filePath, lineNumber, column)
			results <- result
		}

		lineNumber++
		if err == io.EOF {
			break
		}
	}
}

// searchInDirectory walks through a directory and searches each file
func searchInDirectory(rootDir string, searchString string, excludeHidden, caseSensitive bool, results chan<- string, skippedDirs map[string]int) {
	defer wg.Done()

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Exclude hidden files and directories if the -e flag is set
		if excludeHidden && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// If it's a file, search it
		if !info.IsDir() {
			wg.Add(1)
			go searchInFile(path, searchString, caseSensitive, results, skippedDirs)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}

// printResults prints grouped search results
func printResults(results []string) {
	files := make(map[string][]string)

	// Group results by the full file path
	for _, result := range results {
		firstColon := strings.Index(result, ":")
		if firstColon == -1 {
			continue
		}

		secondColon := strings.Index(result[firstColon+1:], ":")
		if secondColon == -1 {
			continue
		}
		secondColon += firstColon + 1

		filePath := result[:secondColon]
		files[filePath] = append(files[filePath], result)
	}

	for filePath, occurrences := range files {
		fmt.Printf("%s (%d):\n", filePath, len(occurrences))
		for _, occurrence := range occurrences {
			fmt.Printf("  - %s\n", occurrence)
		}
	}
}

// printUsage prints usage information for the tool
func printUsage() {
	fmt.Println("grepl version", version)
	fmt.Println("Usage: grepl [-e] [-c] [-s] <search-string>")
	fmt.Println("Flags:")
	fmt.Println("  -e    Do not search subdirectories and hidden files")
	fmt.Println("  -c    Perform case-sensitive search")
	fmt.Println("  -s    Show directories where files have been skipped")
	fmt.Println("\nExample:")
	fmt.Println("  grepl -e -c -s 'search text'")
}

func main() {
	// Define flags for excluding hidden files, case sensitivity, and showing skipped directories
	excludeHidden := flag.Bool("e", false, "Do not search subdirectories and hidden files")
	caseSensitive := flag.Bool("c", false, "Case-sensitive search")
	showSkipped := flag.Bool("s", false, "Show directories where files have been skipped")
	flag.Parse()

	// Print usage help if no arguments are provided
	if len(flag.Args()) != 1 {
		printUsage()
		os.Exit(0)
	}

	searchString := flag.Arg(0)
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}

	// Channel to collect results
	resultsChan := make(chan string)
	var results []string

	// Start a goroutine to search in the root directory concurrently
	wg.Add(1)
	go searchInDirectory(currentDir, searchString, *excludeHidden, *caseSensitive, resultsChan, skippedDirs)

	// Start a goroutine to close the results channel once all searching is done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect the results as they come
	for result := range resultsChan {
		results = append(results, result)
	}

	// Print the results in the desired format
	printResults(results)

	// Print skipped directories if the -s flag is provided
	if *showSkipped && skippedFiles > 0 {
		fmt.Printf("\nSkipped files: %d\n", skippedFiles)
		fmt.Println("Skipped directories:")
		for dir, count := range skippedDirs {
			fmt.Printf("  - %s (%d)\n", dir, count)
		}
	}
}
