package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tliron/commonlog"
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

type schemaModel struct {
	ModelInformation []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Columns     []struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
		} `yaml:"columns"`
	} `yaml:"models"`
}

func (m *schemaModel) ToNode() []Node {
	node := []Node{}
	for _, info := range m.ModelInformation {

		columns := map[string]NodeColumn{}
		for _, column := range info.Columns {
			columns[column.Name] = NodeColumn(column)
		}

		node = append(node, Node{
			Name:        info.Name,
			Description: info.Description,
			Columns:     columns,
		})
	}
	return node
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
	logger := commonlog.GetLoggerf("%s.schema", "settings")
	schemaFiles := []Node{}

	for _, path := range settings.PathSettings.ModelPath {
		modelPath := filepath.Join(settings.RootPath, path)

		err := filepath.Walk(modelPath, func(path string, info fs.FileInfo, error error) error {
			if info.IsDir() {
				return nil
			}

			extension := filepath.Ext(info.Name())

			logger.Infof("Extension : %v", extension)
			if extension != `.yaml` && extension != `.yml` {
				return nil
			}

			fileContent, err := os.ReadFile(path)
			logger.Infof("file : %v", path)
			if err != nil {
				return err
			}

			model := schemaModel{}
			err = yaml.Unmarshal(fileContent, &model)

			logger.Infof("file : %v", model)
			if err != nil {
				return err
			}

			schemaFiles = append(schemaFiles, model.ToNode()...)

			return nil
		})
		if err != nil {
			return nil, err
		}
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
