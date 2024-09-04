package main

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"

	"github.com/alfg/mp4"
)

func collectMediaContent(p string, minDuration int, maxDuration int) ([]MediaItem, error) {
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
			if item.Duration > float64(minDuration) && item.Duration < float64(maxDuration) {
				items = append(items, item)
			}
			idx++
			if idx%500 == 0 {
				fmt.Printf("Processed %d files\n", idx)
			}
		}
		return nil
	})
	return items, err
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

func writePlayList(s any) error {
	outFile, err := os.Create("temp_new.xspf")
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

func buildPlayList(content []MediaItem) *PlayList {
	playList := &PlayList{Xmlns: Xmlns, XmlnsVlc: XmlnsVlc, Version: "1"}
	trackList := &TrackList{}
	tracks := []*Track{}

	for i, media := range content {
		ext := &Extension{Application: ExtensionApplication, Id: i}
		track := &Track{Location: media.AbsPath, Title: media.Name, Duration: math.Round(media.Duration), Ext: *ext}
		tracks = append(tracks, track)
	}
	trackList.Tracks = tracks
	playList.Tl = *trackList
	return playList
}
