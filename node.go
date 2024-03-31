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
	FileUri     string
	Position    protocol.Position
	ProjectName string
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

	positions := parser.GetJinjaPositions(fileString)
	rawPosition := getRawPositionInFile(fileString, params.Position.Line, params.Position.Character)

	// Are we in a jinja block ?
	if !positionWithinRange(rawPosition, positions) {
		return "", nil
	}

	refTags := parser.GetAllRefTags(fileString)
	for _, tag := range refTags {
		if rawPosition >= tag.Range.Start && rawPosition <= tag.Range.End {
			model := fmt.Sprintf("model.%s.%s", params.ProjectName, tag.ModelName)
			return model, nil
		}
	}

	macros := parser.GetMacros(fileString)
	for _, macro := range macros {
		if rawPosition >= macro.Range.Start && rawPosition <= macro.Range.End {
			model := fmt.Sprintf("macro.%s.%s", params.ProjectName, macro.ModelName)
			return model, nil
		}
	}

	return "", nil
}
