package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ProjectSettings struct {
	RootPath     string
	TargetPath   string
	PathSettings pathSettings
}

type pathSettings struct {
	ModelPath []string `yaml:"model-paths"`
	MacroPath []string `yaml:"macro-paths"`
}

type ModelReference struct {
	ModelName string
	Range     Range
}

type Range struct {
	Start int
	End   int
}

type model struct {
	ModelInformation []struct {
		Name        string        `yaml:"name"`
		Description string        `yaml:"description"`
		Columns     []modelColumn `yaml:"columns"`
	} `yaml:"models"`
}

type modelColumn struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func (m *model) ToNode() []Node {
	node := []Node{}
	for _, info := range m.ModelInformation {

		columns := map[string]NodeColumn{}
		for _, column := range info.Columns {
			columns[column.Name] = column.toNodeColumns()
		}

		node = append(node, Node{
			Name:        info.Name,
			Description: info.Description,
			Columns:     columns,
		})
	}
	return node
}

func (m *modelColumn) toNodeColumns() NodeColumn {
	return NodeColumn{
		Name:        m.Name,
		Description: m.Description,
	}
}

func LoadSettings(workspaceFolder string) (ProjectSettings, error) {
	rootDir := strings.ReplaceAll(ROOT_DIR, "file://", "")

	dbtProjectFile := fmt.Sprintf("%s/dbt_project.yml", rootDir)
	targetPath := fmt.Sprintf("%s/target", rootDir)
	fileContent, err := os.ReadFile(dbtProjectFile)
	if err != nil {
		return ProjectSettings{}, err
	}

	settings := pathSettings{}
	err = yaml.Unmarshal([]byte(fileContent), &settings)
	if err != nil {
		return ProjectSettings{}, err
	}

	return ProjectSettings{
		RootPath:     rootDir,
		PathSettings: settings,
		TargetPath:   targetPath,
	}, nil
}

func (settings ProjectSettings) GetSchemaFiles() ([]Node, error) {
	schemaFiles := []Node{}
	for _, path := range settings.PathSettings.ModelPath {
		modelPath := filepath.Join(settings.RootPath, path)

		filepath.Walk(modelPath, func(path string, info fs.FileInfo, error error) error {
			if info.IsDir() {
				return nil
			}

			extenstion := filepath.Ext(path)
			if extenstion != "yaml" || filepath.Ext(path) != "yml" {
				return nil
			}

			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			model := model{}
			err = yaml.Unmarshal(fileContent, &model)
			if err != nil {
				return err
			}

			schemaFiles = append(schemaFiles, model.ToNode()...)

			return nil
		})
	}

	return schemaFiles, nil
}

func (settings ProjectSettings) LoadManifestFile() (Manifest, error) {
	manifestPath := filepath.Join(settings.TargetPath, "manifest.json")
	file, err := os.ReadFile(manifestPath)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	err = json.Unmarshal(file, &manifest)

	return manifest, err
}
