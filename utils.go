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
