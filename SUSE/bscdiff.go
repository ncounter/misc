package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

type searchResult struct {
	line  int
	match []string
	text  string
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <source_file> <target_file>", args[0])
		os.Exit(1)
	}
	// TOOD: Check if files are actually existing.
	searchResults1 := scanFile(args[1])
	searchResults2 := scanFile(args[2])
	missingBscs := findMissingBsc(searchResults1, searchResults2)
	missingBscs = removeDuplicates(missingBscs)
	prettyPrintMissingBscs(searchResults1, missingBscs)
}

// Outputs the missing BSC numbers in a useful format.
func prettyPrintMissingBscs(searchResults1 []searchResult, missingBscs []string) {
	sort.Strings(missingBscs)
	for _, bsc := range missingBscs {
		for _, searchResult := range searchResults1 {
			searchPos := sort.SearchStrings(searchResult.match, bsc)
			if searchPos < len(searchResult.match) && searchResult.match[searchPos] == bsc {
				fmt.Println(fmt.Sprintf("%d: %s -> %s",
					searchResult.line,
					bsc,
					searchResult.text))
			}
		}
	}

}

// Returns a list of BSC numbers, that are missing from the second changelog file.
func findMissingBsc(changelog1 []searchResult, changelog2 []searchResult) []string {
	bscList1 := getBscs(changelog1)
	bscList2 := getBscs(changelog2)
	sort.Strings(bscList1)
	sort.Strings(bscList2)

	var missingBscs []string
	for _, bsc := range bscList1 {
		searchPos := sort.SearchStrings(bscList2, bsc)
		if searchPos < len(bscList2) && bscList2[searchPos] == bsc {
			// found it
		} else {
			missingBscs = append(missingBscs, bsc)
		}
	}
	return missingBscs
}

// Extracts the BSC numbers from the search results and returns them as an array.
func getBscs(res []searchResult) []string {
	var bsc []string
	for _, v := range res {
		for _, value := range v.match {
			bsc = append(bsc, value)
		}
	}
	return bsc
}

// Scans the file for BSC numbers and returns the search results.
func scanFile(pathToFile string) []searchResult {
	var re, _ = regexp.Compile(`bsc#\d*`)
	lines, err := scanLines(pathToFile)
	if err != nil {
		panic(err)
	}
	var searchResults []searchResult
	for i, line := range lines {
		results := re.FindAllString(line, -1)
		if len(results) > 0 {
			res := searchResult{
				line:  i + 1,
				match: results,
				text:  line}
			searchResults = append(searchResults, res)
		}
	}
	return searchResults
}

// Returns the given file as an array of lines.
func scanLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

// Removed duplicates form an array.
func removeDuplicates(s []string) []string {
	m := make(map[string]bool)
	for _, item := range s {
		if _, ok := m[item]; ok {
			// duplicate item
		} else {
			m[item] = true
		}
	}
	var result []string
	for item := range m {
		result = append(result, item)
	}
	return result
}