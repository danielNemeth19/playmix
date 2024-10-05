package main

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	// "syscall"
	// "time"

	"github.com/alfg/mp4"
)

// TODO: let's sanitize track title by cutting the vlc record prefix
type MediaItem struct {
	AbsPath  string
	Dir      string
	Name     string
	Id       int
	Duration float64
}

// TODO: This dirName could be used writing a proper title
func (m MediaItem) getRelativeDir(rootParts []string) string {
	fileParts := getPathParts(m.AbsPath)

	if len(fileParts) == len(rootParts) {
		return filepath.Base(filepath.Dir(m.AbsPath))
	} else {
		res := fileParts[len(rootParts):]
		return filepath.Join(res...)
	}
}

type Summarizer struct {
	dBucket       DurationBucket
	ratio         int
	totalDuration float64
	totalScanned  int
	totalSelected int
}

func (s Summarizer) getRealRatio() float64 {
	return float64(s.totalSelected) / float64(s.totalScanned) * 100
}

func (s Summarizer) getData() {
	fmt.Printf("Total scanned: %d\n", s.totalScanned)
	fmt.Printf("Duration distribution:\n")
	s.dBucket.summarize()
	fmt.Printf("Total duration is: %f sec -- (%f) minutes\n", s.totalDuration, s.totalDuration/60)
	fmt.Printf("Total selected: %d -- required ratio: %d -- got: %.2f%%\n", s.totalSelected, s.ratio, s.getRealRatio())
}

type DurationBucket struct {
	Dur0_5     int
	Dur5_10    int
	Dur10_30   int
	Dur30_60   int
	Dur60_180  int
	Dur180_240 int
	DurOver240 int
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
	case duration <= 180:
		d.Dur60_180++
	case duration <= 240:
		d.Dur180_240++
	default:
		d.DurOver240++
	}
}

func (d *DurationBucket) summarize() {
	fmt.Printf("Bucket <5 seconds: %d\n", d.Dur0_5)
	fmt.Printf("Bucket 5-10 seconds: %d\n", d.Dur5_10)
	fmt.Printf("Bucket 10-30 seconds: %d\n", d.Dur10_30)
	fmt.Printf("Bucket 30-60 seconds: %d\n", d.Dur30_60)
	fmt.Printf("Bucket 60-180 seconds: %d\n", d.Dur60_180)
	fmt.Printf("Bucket 180-240 seconds: %d\n", d.Dur180_240)
	fmt.Printf("Bucket 240< seconds: %d\n", d.DurOver240)
}

func selector(ratio int) bool {
	n := rand.Intn(100)
	if n < ratio {
		return true
	}
	return false
}

func toSkip(name string, skip []string) bool {
	for _, f := range skip {
		if f == name {
			return true
		}
	}
	return false
}

// TODO: maybe use the path separator utility here too?
func isIncluded(root string, path string, include []string) bool {
	if len(include) == 0 {
		return true
	}
	parts := strings.Split(path, root)
	for _, folder := range strings.Split(parts[1], string(filepath.Separator)) {
		for _, f := range include {
			if folder == f {
				return true
			}
		}
	}
	return false
}

func collectMediaContent(p string, params Params) ([]MediaItem, Summarizer, error) {
	var items []MediaItem
	rootParts := getPathParts(p)
	summary := Summarizer{
		ratio:         params.ratio,
		totalScanned:  0,
		totalSelected: 0,
		totalDuration: 0,
		dBucket:       DurationBucket{},
	}
	idx := 0
	err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && toSkip(d.Name(), params.skipF) {
			return filepath.SkipDir
		}
		if !d.IsDir() && isMediaFile(filepath.Ext(d.Name())) && isIncluded(p, path, params.includeF) {
			// file, _ := d.Info()
			// ctime := file.Sys().(*syscall.Stat_t).Ctim
			// cTime := time.Unix(ctime.Sec, ctime.Nsec)
			// fmt.Printf("Change time: %v\n", cTime)
			// mTime := file.ModTime()
			// fmt.Printf("Modified time: %v\n", mTime)
			// fmt.Printf("%s -- %v -- %v\n", file.Name(), cTime, mTime)

			if selector(params.ratio) {
				duration, err := getDuration(path)
				summary.dBucket.allocate(duration)
				if err != nil {
					return err
				}
				if duration > float64(params.minDuration) && duration < float64(params.maxDuration) {
					item := MediaItem{Id: idx, AbsPath: path, Name: d.Name(), Duration: duration}
					item.Dir = item.getRelativeDir(rootParts)
					items = append(items, item)
					summary.totalDuration += duration
					summary.totalSelected++
				}
			}
			summary.totalScanned++
			if summary.totalScanned%500 == 0 {
				fmt.Printf("Processed %d files\n", summary.totalScanned)
			}
		}
		return nil
	})
	return items, summary, err
}

func getDuration(p string) (duration float64, err error) {
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

	duration = rawDuration / timeScale
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
