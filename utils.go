package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var fileRegex = regexp.MustCompile(`file:`)

func getModelNameFromFilePath(filePath string) string {
	file := filepath.Base(filePath)
	file = file[:strings.Index(file, ".")]
	return file
}

func getRawPositionInFile(content string, line, character uint32) int {
	// where are we in the file
	position := 0
	fileLines := strings.Split(content, "\n")
	for i := uint32(0); i < line; i++ {
		position += len(fileLines[i])
		fmt.Println("line", i)
	}
	position += int(character)
	return position
}

func positionWithinRange(rawPosition int, ranges []Range) bool {
	for _, r := range ranges {
		if rawPosition >= r.Start && rawPosition <= r.End {
			return true
		}
	}
	return false
}

func ReadFileUri(fileUri string) ([]byte, error) {
	u, _ := url.ParseRequestURI(fileUri)
	return os.ReadFile(u.Path)
}
