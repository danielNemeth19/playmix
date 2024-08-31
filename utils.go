package main

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alfg/mp4"
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

func getDuration(p string) (float64, error) {
	// log.Println("Checking ", p)
	file, err := os.Open(p)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	mp4, err := mp4.OpenFromReader(file, info.Size())
	if err != nil {
		return 0, err
	}
	if mp4.Moov == nil {
		return 0, fmt.Errorf("Moov box not found for %s. Is this mp4?", p)
	}
	rawDuration := float64(mp4.Moov.Mvhd.Duration)
	timeScale := float64(mp4.Moov.Mvhd.Timescale)

	duration := rawDuration / timeScale
	return duration, nil
}

func collectMediaContent(p string) ([]MediaItem, error) {
	var items []MediaItem
	idx := 0
	err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && isMediaFile(filepath.Ext(d.Name())) {
			duration, err := getDuration(path)
			if err != nil {
				return err
			}
			item := MediaItem{Id: idx, AbsPath: path, Name: d.Name(), Duration: duration}
			items = append(items, item)
			idx++
            if idx % 500 == 0 {
                fmt.Printf("Processed %d files\n", idx)
            }
		}
		return nil
	})
	return items, err
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
