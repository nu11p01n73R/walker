package walker

import (
	"os"
	"strings"
	"testing"
)

var TEST_FILES = [6]string{
	"one/two/three/four",
	"one/two/three/five",
	"one/two/six/seven",
	"one/two/six/eight",
	"one/two/nine",
	"one/ten"}

func getBase(path string) string {
	index := strings.LastIndex(path, "/")
	return path[:index]
}

func createFile(fileName string) error {
	path := getBase(fileName)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(fileName, os.O_CREATE, 0755)
	defer file.Close()
	return err
}

func createFiles() error {
	for _, file := range TEST_FILES {
		err := createFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func verifyFiles(walked []string) bool {
	if len(walked) != len(TEST_FILES) {
		return false
	}

	hash := map[string]bool{}
	for _, file := range walked {
		hash[file] = true
	}

	for _, file := range TEST_FILES {
		_, ok := hash[file]
		if !ok {
			return false
		}
	}

	return true
}

func TestWalker(t *testing.T) {
	err := createFiles()
	if err != nil {
		t.Error(err.Error())
	}

	files, _ := Walk("one", func(files []string) []string {
		return files
	}, []string{`git`, `sw`})

	if !verifyFiles(files) {
		t.Error("Walked files not same as created files")
	}

	// Cleanup
	err = os.RemoveAll("one")
	if err != nil {
		t.Error(err.Error())
	}
}
