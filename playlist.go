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

type DurationBucket struct {
	Dur0_5    int
	Dur5_10   int
	Dur10_30  int
	Dur30_60  int
	DurOver60 int
}

func (d *DurationBucket) allocate(duration float64) {
	switch {
	case duration < 5:
		d.Dur0_5++
	case duration <= 10:
		d.Dur5_10++
	case duration <= 30:
		d.Dur10_30++
	case duration <= 60:
		d.Dur30_60++
	default:
		d.DurOver60++
	}
}

func (d DurationBucket) summarize() {
	fmt.Printf("Bucket <5 seconds: %d\n", d.Dur0_5)
	fmt.Printf("Bucket 5-10 seconds: %d\n", d.Dur5_10)
	fmt.Printf("Bucket 10-30 seconds: %d\n", d.Dur10_30)
	fmt.Printf("Bucket 30-60 seconds: %d\n", d.Dur30_60)
	fmt.Printf("Bucket 60< seconds: %d\n", d.DurOver60)
}

func collectMediaContent(p string, minDuration int, maxDuration int) ([]MediaItem, error) {
	var items []MediaItem 
    var totalDuration float64
	durationMap := &DurationBucket{Dur0_5: 0, Dur5_10: 0, Dur10_30: 0, Dur30_60: 0, DurOver60: 0}
	idx := 0
	err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && isMediaFile(filepath.Ext(d.Name())) {
			duration, err := getDuration(path)
			durationMap.allocate(duration)
			if err != nil {
				return err
			}
			if duration > float64(minDuration) && duration < float64(maxDuration) {
				item := MediaItem{Id: idx, AbsPath: path, Name: d.Name(), Duration: duration}
				items = append(items, item)
                totalDuration += duration
			}
			idx++
			if idx%500 == 0 {
				fmt.Printf("Processed %d files\n", idx)
			}
		}
		return nil
	})
	durationMap.summarize()
    fmt.Printf("Total duration is: %f sec -- (%f) minutes\n", totalDuration, totalDuration / 60)
	return items, err
}

func getDuration(p string) (float64, error) {
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
