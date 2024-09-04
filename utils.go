package main

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

func getPath() (string, error) {
	path := os.Getenv("MEDIA_SOURCE")
	if path == "" {
		return "", fmt.Errorf("MEDIA_SOURCE environment variable not set")
	}
	return path, nil
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func dumpConsole(s any) {
	out, _ := xml.MarshalIndent(s, " ", "  ")
	fmt.Println(xml.Header + string(out))
}

func isMediaFile(ext string) bool {
	for _, v := range mediaExtensions {
		if v == ext {
			return true
		}
	}
	return false
}

func collectExtensions(p string) ([]string, error) {
	extensions := []string{}
	seen := map[string]bool{}

	err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			extension := filepath.Ext(path)
			if !seen[extension] {
				fmt.Printf("Seen first: %s -- %s, %v\n", extension, d.Name(), d.IsDir())
				seen[extension] = true
				extensions = append(extensions, extension)
			}
		}
		return nil
	})
	return extensions, err
}

func _collectExtensions(p string, extensions *[]string, seen *map[string]bool) (*[]string, error) {
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}
	for _, obj := range files {
		if obj.IsDir() {
			extensions, _ = _collectExtensions(p+"/"+obj.Name(), extensions, seen)
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
	extensions, err := _collectExtensions(p, extensions, seen)
	return extensions, err
}
