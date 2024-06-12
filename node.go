package main

import (
	"fmt"

	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Manifest struct {
	Nodes    map[string]Node  `json:"nodes"`
	Macros   map[string]Macro `json:"macros"`
	Metadata Metadata         `json:"metadata"`
}

type Metadata struct {
	ProjectName string `json:"project_name"`
}

type Node struct {
	Columns      map[string]NodeColumn
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	OriginalPath string  `json:"original_file_path"`
	RawCode      string  `json:"raw_code"`
	Depends      Depends `json:"depends_on"`
}

type Macro struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	OriginalPath string `json:"original_file_path"`
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
	ProjectName string
	Manifest    Manifest
	Position    protocol.Position
}

type DefinitionResponse struct {
	FileName string
}

func (n Node) GetDefinition(params DefinitionRequest) (DefinitionResponse, error) {
	logger := commonlog.GetLogger("node.GetDefinition")
	parser := NewJinjaParser()

	fileContent, err := ReadFileUri(params.FileUri)
	if err != nil {
		logger.Infof("couldn't read file %v", err)
		return DefinitionResponse{}, err
	}

	fileString := string(fileContent)
	if parser.HasJinjaBlocks(fileString) {
		positions := parser.GetJinjaPositions(fileString)
		rawPosition := getRawPositionInFile(fileString, params.Position.Line, params.Position.Character)

		// Are we in a jinja block ?
		if positionWithinRange(rawPosition, positions) {
			logger.Info("We are within a jinja block")
			return n.getJinjaDefinition(params, rawPosition, fileString, parser)
		} else {
			logger.Info("Not in a jinja block")
		}
	} else {
		logger.Info("No jinja block found")
	}

	// handle sql definition

	return DefinitionResponse{}, nil
}

func (n Node) getJinjaDefinition(params DefinitionRequest, rawPosition int, content string, parser JinjaParser) (DefinitionResponse, error) {
	logger := commonlog.GetLogger("lsp.getJinjaDefinition")

	refTags := parser.GetAllRefTags(content)
	for _, tag := range refTags {
		if rawPosition >= tag.Range.Start && rawPosition <= tag.Range.End {
			model := fmt.Sprintf("model.%s.%s", params.ProjectName, tag.ModelName)
			node, ok := params.Manifest.Nodes[model]

			logger.Infof("looking for model %v", model)
			if !ok {
				return DefinitionResponse{}, nil
			}
			return DefinitionResponse{FileName: node.OriginalPath}, nil
		}
	}
	logger.Info("not within ref tag")

	macros := parser.GetMacros(content)
	logger.Infof("could not find a ref tag trying macro %v", macros)
	for _, macro := range macros {
		if rawPosition >= macro.Range.Start && rawPosition <= macro.Range.End {
			model := fmt.Sprintf("macro.%s.%s", params.ProjectName, macro.ModelName)
			node, ok := params.Manifest.Macros[model]

			logger.Infof("looking for macro %v", model)
			if !ok {
				return DefinitionResponse{}, nil
			}
			return DefinitionResponse{FileName: node.OriginalPath}, nil
		}
	}

	logger.Info("not withing macro")

	return DefinitionResponse{}, nil
}
