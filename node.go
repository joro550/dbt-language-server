package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/tliron/commonlog"
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

func (n Node) GetDefinition(params *protocol.DefinitionParams) error {
	logger := commonlog.GetLogger("node.GetDefinition")
	parser := NewJinjaParser()

	fileContent, err := os.ReadFile(strings.ReplaceAll(params.TextDocument.URI, "file://", ""))
	if err != nil {

		logger.Infof("couldn't read fire %v", err)
		return err
	}

	fileString := string(fileContent)
	if !parser.HasJinjaBlocks(fileString) {
		logger.Info("doesn't have jinja blocks")
		return nil
	}

	// where are we in the file
	position := 0
	fileLines := strings.Split(fileString, "\n")
	for i := uint32(0); i < params.Position.Line; i++ {
		position += len(fileLines[i])
	}
	position += int(params.Position.Character)
	logger.Infof("calculated position %d", position)

	refTags := parser.GetAllRefTags(fileString)
	logger.Infof("reftags %d", refTags)

	return nil
}
