package main

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getPath() (string, error) {
	rootPath := os.Getenv("MEDIA_SOURCE")
	if rootPath == "" {
		return "", fmt.Errorf("MEDIA_SOURCE environment variable not set")
	}
	if !strings.HasSuffix(rootPath, string(filepath.Separator)) {
		rootPath += string(filepath.Separator)
		log.Printf("Root path got normalized by adding path separator (%s)\n", string(filepath.Separator))
	}
	return rootPath, nil
}

func getPathParts(p string) []string {
	dir, _ := filepath.Split(p)
	parts := strings.Split(dir, string(filepath.Separator))
	return parts[1 : len(parts)-1]
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
