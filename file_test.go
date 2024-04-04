package main

import (
	"fmt"
	"net/url"
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

func TestFileURI(t *testing.T) {
	for _, path := range []string{
		"file:///path/to/file.json",
		"file:///c:/WINDOWS/clock.json",
		"file://localhost/path/to/file.json",
		"file://localhost/c:/WINDOWS/clock.avi",

		// A case that you probably don't need to handle given the rarity,
		// but is a known legacy win32 issue when translating \\remotehost\share\dir\file.txt
		"file:////remotehost/share/dir/file.txt",
	} {
		u, _ := url.ParseRequestURI(path)
		fmt.Printf("url:%v\nscheme:%v host:%v Path:%v\n\n", u, u.Scheme, u.Host, u.Path)
	}
}
