package walker

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// Function which formats and/or filters
// the file list.
type Filter func(files []string) []string

// Recursively walks the directories from
// dir.
// Return
//	array of strings of files
func getFiles(dir string, ignoreRegex *regexp.Regexp) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	var output []string
	for _, file := range files {
		if ignoreRegex != nil && ignoreRegex.MatchString(file.Name()) {
			continue
		}

		fullName := fmt.Sprintf("%s/%s", dir, file.Name())
		if file.IsDir() {
			output = append(output, getFiles(fullName, ignoreRegex)...)
		} else {
			output = append(output, fullName)
		}
	}
	return output
}

// Generate regex pattern from the list of file
// to be ignored
// Return
//	Compiled regexp
func buildPattern(ignore []string) *regexp.Regexp {
	pattern := strings.Join(ignore, `|`)

	return regexp.MustCompile(pattern)
}

// Returns a list of files from the root directory.
// Walk iterates through all the subdirectories and return
// a list of regular files.
//
// For each of the directory at the root, a new go routine is
// started, and they are iterated recursively and parallelly.
//
// formatFunc is called at the end of the walk to filter and/or format
// the files.
//
// ignore, list of file patterns to be ignored from the parent directory.
// Accepts regular expressions. If nothing to be ignored, has to be
// provided as empty array.
//
// Returns an array of string, files and
// error if some error occured.
func Walk(root string, formatFunc Filter, ignore []string) ([]string, error) {
	var ignoreRegex *regexp.Regexp
	if len(ignore) > 0 {
		ignoreRegex = buildPattern(ignore)
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return []string{}, err
	}

	var output []string
	var workerCount int
	queue := make(chan []string)
	for _, file := range files {
		if ignoreRegex != nil && ignoreRegex.MatchString(file.Name()) {
			continue
		}

		fullName := fmt.Sprintf("%s/%s", root, file.Name())
		if !file.IsDir() {
			output = append(output, fullName)
		} else {
			go func() {
				files := getFiles(fullName, ignoreRegex)
				queue <- files
			}()
			workerCount++
		}
	}

	for ; workerCount > 0; workerCount-- {
		files := <-queue
		output = append(output, files...)
	}

	return formatFunc(output), err
}
