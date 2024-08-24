package main

import (
	"fmt"
	"encoding/xml"
	"os"
	"path/filepath"
)

func getPath() (string, error) {
	path := os.Getenv("MEDIA_SOURCE")
	if path == "" {
		return "", fmt.Errorf("MEDIA_SOURCE environment variable not set")
	}
	return path, nil
}

func dumpConsole(s any) {
	out, _ := xml.MarshalIndent(s, " ", "  ")
	fmt.Println(xml.Header + string(out))
}

func collectExtensions(p string, extensions *[]string, seen *map[string]bool) (*[]string, error) {
	fmt.Println("Checking folder: ", p)
	files, err := os.ReadDir(p)
	if err != nil {
        return nil, err
	}
	for _, obj := range files {
		if obj.IsDir() {
			extensions, _ = collectExtensions(p+"/"+obj.Name(), extensions, seen)
		} else {
			extension := filepath.Ext(obj.Name())
			if !(*seen)[extension] {
				fmt.Printf("Seen first: %s -- %s, %v\n", extension, obj.Name(), obj.IsDir())
				(*seen)[extension] = true
				*extensions = append(*extensions, extension)
			}
		}
	}
	return extensions, nil
}

func getExtensions(p string) (*[]string, error) {
    extensions := &[]string{}
    seen := &map[string]bool{}
    extensions, err := collectExtensions(p, extensions, seen)
    return extensions, err
}
