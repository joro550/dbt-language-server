package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"
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
	Description  string  `json:"description"`
	OriginalPath string  `json:"original_file_path"`
	RawCode      string  `json:"raw_code"`
	Depends      Depends `json:"depends_on"`
	Columns      map[string]NodeColumn
}
type NodeColumn struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Depends struct {
	Nodes []string `json:"nodes"`
}

type DefinitionRequest struct {
	FileUri  string
	Position protocol.Position
}

func (n Node) GetDefinition(params DefinitionRequest) (string, error) {
	logger := commonlog.GetLogger("node.GetDefinition")
	parser := NewJinjaParser()

	fileContent, err := os.ReadFile(strings.ReplaceAll(params.FileUri, "file://", ""))
	if err != nil {

		logger.Infof("couldn't read file %v", err)
		return "", err
	}

	fileString := string(fileContent)
	if !parser.HasJinjaBlocks(fileString) {
		logger.Info("doesn't have jinja blocks")
		return "", nil
	}

	refTags := parser.GetAllRefTags(fileString)
	position := getRawPositionInFile(fileString, params.Position.Line, params.Position.Character)

	for _, tag := range refTags {
		if position >= tag.Range.Start && position <= tag.Range.End {
			return tag.ModelName, nil
		}
	}
	return "", nil
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
