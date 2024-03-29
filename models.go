package main

import (
	"encoding/json"
	"fmt"
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
