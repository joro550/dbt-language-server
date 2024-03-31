package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilePath(t *testing.T) {
	filePath := `file://root/dev/dbt/file.sql`

	base := filepath.Base(filePath)
	directory := filepath.Dir(filePath)

	index := strings.Index(base, ".")
	fmt.Println(index)
	fmt.Println(directory)

	base = base[:index]

	if base != "file" {
		t.Fatalf("got %v", base)
	}
}
