package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func readInOptFile(fsys fs.FS, fn string) ([]byte, error) {
	file, err := fsys.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("File cannot be opened: %s", fn)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("File cannot be read: %s", fn)
	}
	return data, nil
}
func getPathParts(p string) []string {
	dir := filepath.Dir(p)
	trimmedDir := strings.TrimPrefix(dir, string(filepath.Separator))
	parts := strings.Split(trimmedDir, string(filepath.Separator))
	return parts
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

func getUrlEncodedPath(path string) string {
	dir, fn := filepath.Split(path)
	location := "file:///" + dir + url.PathEscape(fn)
	return location
}

func collectExtensions(fsys fs.FS) ([]string, error) {
	extensions := []string{}
	seen := map[string]bool{}

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			extension := filepath.Ext(path)
			if !seen[extension] {
				seen[extension] = true
				extensions = append(extensions, extension)
			}
		}
		return nil
	})
	return extensions, err
}

func createFile(fn string) (*os.File, error) {
	outFile, err := os.Create(fn)
	if err != nil {
		return nil, fmt.Errorf("Error creating file: %w\n", err)
	}
	return outFile, nil
}
