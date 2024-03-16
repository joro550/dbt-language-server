package main

import (
	"regexp"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	jinjaNode, _ = regexp.Compile("{{[\\s\\S]*}}")
	name, _      = regexp.Compile("(?:'|\")[\\S]*(?:'|\")")
)

type Manifest struct {
	Nodes    map[string]Node `json:"nodes"`
	Metadata Metadata        `json:"metadata"`
}

type Metadata struct {
	ProjectName string `json:"project_name"`
}

type Node struct {
	Name         string  `json:"name"`
	OriginalPath string  `json:"original_file_path"`
	Depends      Depends `json:"depends_on"`
	RawCode      string  `json:"raw_code"`
}

type Depends struct {
	Nodes []string `json:"nodes"`
}

func (n Node) DoThing(currentPosition protocol.Position) (bool, string) {
	lines := strings.Split(n.RawCode, "\n")

	currentLine := lines[currentPosition.Line]
	currentLineAsBytes := []byte(currentLine)
	match := jinjaNode.Match(currentLineAsBytes)

	if !match {
		return false, currentLine
	}

	indicies := name.FindAllIndex(currentLineAsBytes, -1)
	if indicies == nil {
		return false, currentLine
	}

	for _, index := range indicies {
		start := uint32(index[0])
		end := uint32(index[1])

		if currentPosition.Character >= start && currentPosition.Character <= end {
			reference := currentLine[start+1 : end-1]
			return true, reference
		}
	}

	return false, currentLine
}

func (n Node) DoThing2(code string, currentPosition protocol.Position) (bool, string) {
	lines := strings.Split(code, "\n")

	currentLine := lines[currentPosition.Line]
	currentLineAsBytes := []byte(currentLine)
	match := jinjaNode.Match(currentLineAsBytes)

	if !match {
		return false, currentLine
	}

	indicies := name.FindAllIndex(currentLineAsBytes, -1)
	if indicies == nil {
		return false, currentLine
	}

	for _, index := range indicies {
		start := uint32(index[0])
		end := uint32(index[1])

		if currentPosition.Character >= start && currentPosition.Character <= end {
			reference := currentLine[start+1 : end-1]
			return true, reference
		}
	}

	return false, currentLine
}
