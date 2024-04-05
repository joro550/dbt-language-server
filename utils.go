package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

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
	u, _ := CleanUri(fileUri)
	return os.ReadFile(u)
}

func ReadFileUri2(path, file string) ([]byte, error) {
	u, _ := CleanUri(path)
	return os.ReadFile(filepath.Join(u, file))
}

func CleanUri(fileUri string) (string, error) {
	driveRegex := regexp.MustCompile(`[a-zA-Z]:\/\/`)
	cleanedUri, err := url.ParseRequestURI(fileUri)

	var cleanedPath string
	// This is probably because this is not a uri
	if err != nil {
		cleanedPath = fileUri
	} else {
		cleanedPath = cleanedUri.Path
	}

	// this is basically a "are we in windows" check
	if driveRegex.MatchString(fileUri) {
		if cleanedPath[0] == '/' || cleanedPath[0] == '\\' {
			cleanedPath = cleanedPath[1:]
		}
	}

	return cleanedPath, nil
}
