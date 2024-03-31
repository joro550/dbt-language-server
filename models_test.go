package main

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSchemaFile(t *testing.T) {
	fileContent, _ := os.ReadFile("./tests/schema.yml")
	model := schemaModel{}
	err := yaml.Unmarshal(fileContent, &model)
	if err != nil {
		t.Errorf("error %v", err)
	}
	modelCount := len(model.ModelInformation)

	if modelCount != 2 {
		t.Errorf("expected 2 but got %v", modelCount)
	}
}
