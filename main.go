package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/alfg/mp4"
)

const (
	ExtensionApplication = "http://www.videolan.org/vlc/playlist/0"
	Xmlns                = "http://xspf.org/ns/0/"
	XmlnsVlc             = "http://www.videolan.org/vlc/playlist/ns/0/"
)

var mediaExtensions = []string{".mp4", ".mkv", ".avi", ".flv", ".mpeg"}

func buildTrackList(p string) *TrackList {
	files, err := os.ReadDir(p)
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	trackList := &TrackList{}
	var tracks []*Track
	for i, v := range files {
		ext := &Extension{Application: ExtensionApplication, Id: i}
		location := "file://" + p + v.Name()
		track := &Track{Location: location, Ext: *ext}
		tracks = append(tracks, track)
	}
	trackList.Tracks = tracks
	return trackList
}

func writePlayList(s any) error {
	outFile, err := os.Create("temp.xspf")
	if err != nil {
		return fmt.Errorf("Error creating file: %w\n", err)
	}
	defer outFile.Close()

	_, err = outFile.WriteString(xml.Header)
	if err != nil {
		return fmt.Errorf("Error writing header: %w\n", err)
	}

	encoder := xml.NewEncoder(outFile)
	encoder.Indent("", "\t")
	err = encoder.Encode(&s)
	if err != nil {
		return fmt.Errorf("Error in encoding xml: %w\n", err)
	}
	return nil
}

func getDuration(p string) (float64, error) {
    log.Println("Checking ", p)
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
        return 0, fmt.Errorf("Moov box not found.. is this mp4?")
    }
    rawDuration := float64(mp4.Moov.Mvhd.Duration)
    timeScale := float64(mp4.Moov.Mvhd.Timescale)

    duration := rawDuration / timeScale
	return duration, nil
}

func recGgetFolderContent(p string) ([]MediaItem, error) {
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
			item := MediaItem{Id: idx, Name: d.Name(), Duration: duration}
			items = append(items, item)
            idx++
        }
        return nil
    })
    return items, err
}

func main() {
	path, err := getPath()
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	log.Printf("Path to be used: %s\n", path)

	extFlag := flag.Bool("ext", false, "If specified, collects unique file extensions")
	flag.Parse()

	if *extFlag {
		extensions, err := getExtensions(path)
		if err != nil {
			log.Fatalf("Error during extension collection: %s\n", err)
		}
		fmt.Printf("Extensions: %v\n", *extensions)
	}

	content, err := recGgetFolderContent(path)
	if err != nil {
		log.Fatalf("Error during getting files: %s\n", err)
	}
    // fmt.Printf("len content: %d\n", len(content))
    for _, item := range content {
        fmt.Printf("Id: %d -- Name: %s -- %f sec\n", item.Id, item.Name, item.Duration)
    }
}
