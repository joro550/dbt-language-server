package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/tliron/commonlog"
	"gopkg.in/yaml.v3"
)

type ProjectSettings struct {
	Name         string
	RootPath     string
	TargetPath   string
	PathSettings pathSettings
}

type pathSettings struct {
	Name      string   `yaml:"name"`
	ModelPath []string `yaml:"model-paths"`
	MacroPath []string `yaml:"macro-paths"`
}

type ModelReference struct {
	ModelName string
	Range     Range
}

type MacroReference struct {
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
	cleanedWorkspaceUri, err := CleanUri(workspaceFolder)
	if err != nil {
		return ProjectSettings{}, err
	}

	fileContent, err := ReadFileUri2(cleanedWorkspaceUri, "dbt_project.yml")
	if err != nil {
		return ProjectSettings{}, err
	}

	settings := pathSettings{}
	err = yaml.Unmarshal([]byte(fileContent), &settings)
	if err != nil {
		return ProjectSettings{}, err
	}
	return ProjectSettings{
		Name:         settings.Name,
		RootPath:     cleanedWorkspaceUri,
		PathSettings: settings,
		TargetPath:   filepath.Join(cleanedWorkspaceUri, "target"),
	}, nil
}

func (ps ProjectSettings) GetRootDirectory() string {
	return ps.RootPath
}

func (settings ProjectSettings) GetSchemaFiles() ([]Node, error) {
	logger := commonlog.GetLoggerf("%s.schema", "settings")
	schemaFiles := []Node{}

	for _, path := range settings.PathSettings.ModelPath {
		modelPath := filepath.Join(settings.GetRootDirectory(), path)

		err := filepath.Walk(modelPath, func(path string, info fs.FileInfo, error error) error {
			if info.IsDir() {
				return nil
			}

			extension := filepath.Ext(info.Name())

			logger.Infof("Extension : %v", extension)
			if extension != `.yaml` && extension != `.yml` {
				return nil
			}

			fileContent, err := ReadFileUri(path)
			logger.Infof("file : %v", path)
			if err != nil {

				logger.Infof("Could not read file: %v", err)
				return err
			}

			model := schemaModel{}
			err = yaml.Unmarshal(fileContent, &model)

			logger.Infof("file : %v", model)
			if err != nil {
				logger.Infof("Could not parse yaml file", err)
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

func (settings ProjectSettings) PredictManifestFile(projectName string) (Manifest, error) {
	logger := commonlog.GetLogger("models.PredictManifestFile")

	manifest := Manifest{
		Nodes:    map[string]Node{},
		Macros:   map[string]Macro{},
		Metadata: Metadata{ProjectName: projectName},
	}

	for _, path := range settings.PathSettings.ModelPath {
		modelPath := filepath.Join(settings.GetRootDirectory(), path)

		filepath.Walk(modelPath, func(path string, info fs.FileInfo, error error) error {
			if info.IsDir() {
				return nil
			}

			extension := filepath.Ext(info.Name())
			if extension != `.sql` {
				return nil
			}

			fileContent, err := ReadFileUri(path)
			logger.Infof("file : %v", path)
			if err != nil {
				logger.Infof("Could not read file: %v", err)
				return err
			}

			fileName := strings.ReplaceAll(info.Name(), ".sql", "")
			node := Node{
				Name:    fileName,
				RawCode: string(fileContent),
				Columns: map[string]NodeColumn{},
			}

			key := fmt.Sprintf("model.%v.%v", projectName, fileName)
			manifest.Nodes[key] = node

			return nil
		})
	}

	return manifest, nil
}

func (settings ProjectSettings) LoadManifestFile() (Manifest, error) {
	manifestPath := filepath.Join(settings.TargetPath, "manifest.json")
	file, err := ReadFileUri2(manifestPath, "manifest.json")
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	err = json.Unmarshal(file, &manifest)

	return manifest, err
}
