package main

import (
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tliron/commonlog"
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
		position += len(fileLines[i]) + 1
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
	logger := commonlog.GetLogger("utils.ReadFileUri")
	logger.Infof("fileUri: %s", fileUri)
	u, err := CleanUri(fileUri)
	logger.Infof("cleaned uri: %s", u)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(u)
}

func ReadFileUri2(path, file string) ([]byte, error) {
	u, err := CleanUri(path)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filepath.Join(u, file))
}

func CleanUri(fileUri string) (string, error) {
	driveRegex := regexp.MustCompile(`[a-zA-Z]:\/\/`)
	cleanedUri, err := url.ParseRequestURI(fileUri)

	var cleanedPath string
	// This is probably because this is not a uri
	if err != nil || cleanedUri.Path == "" {
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
