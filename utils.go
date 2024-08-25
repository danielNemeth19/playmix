package main

import (
	"fmt"
	"io/fs"
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

func isMediaFile(ext string) bool {
	for _, v := range mediaExtensions {
		if v == ext {
			return true
		}
	}
	return false
}

func getFolderContent(p string) ([]MediaItem, error) {
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}
	var content []MediaItem
	for i, file := range files {
		extension := filepath.Ext(file.Name())
        // could have a condition check for item being a file..
		if isMediaFile(extension) {
			duration, err := getDuration(p + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			item := MediaItem{Id: i, Name: file.Name(), Duration: duration}
			content = append(content, item)
		} else {
            fmt.Println("not a video file:", file.Name())
        }
	}
	return content, nil
}
func walkCollectExtensions(p string) (*[]string, error) {
    extensions := &[]string{}
    seen := &map[string]bool{}

    err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {
            extension := filepath.Ext(path)
			if !(*seen)[extension] {
                fmt.Printf("Seen first: %s -- %s, %v\n", extension, d.Name(), d.IsDir())
				(*seen)[extension] = true
				*extensions = append(*extensions, extension)
			}
        }
        return nil
    })
    return extensions, err
}

func collectExtensions(p string, extensions *[]string, seen *map[string]bool) (*[]string, error) {
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
