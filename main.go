package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	ExtensionApplication = "http://www.videolan.org/vlc/playlist/0"
	Xmlns                = "http://xspf.org/ns/0/"
	XmlnsVlc             = "http://www.videolan.org/vlc/playlist/ns/0/"
)

var mediaExtensions = []string{".mp4", ".mkv", ".avi", ".flv", ".mpeg"}

func buildPlayList(content []MediaItem) *PlayList {
    playList := &PlayList{Xmlns: Xmlns, XmlnsVlc: XmlnsVlc, Version: "1"}
	trackList := &TrackList{}
    tracks := []*Track{}

    for i, media := range content {
        fmt.Printf("index %d - name %s  -- duration: %f\n", i, media.Name, media.Duration)
        ext := &Extension{Application: ExtensionApplication, Id: i}
        track := &Track{Location: media.AbsPath, Title: media.Name, Duration: media.Duration, Ext: *ext}
        tracks = append(tracks, track)
    }
    trackList.Tracks = tracks
    playList.Tl = *trackList
    return playList
}

func buildTrackListBak(p string) *TrackList {
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

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
    defer TimeTrack(time.Now(), "main")
	extFlag := flag.Bool("ext", false, "If specified, collects unique file extensions")
	flag.Parse()

	path, err := getPath()
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	log.Printf("Path to be used: %s\n", path)

	if *extFlag {
		extensions, err := collectExtensions(path)
		if err != nil {
			log.Fatalf("Error during extension collection: %s\n", err)
		}
		fmt.Printf("Extensions: %v\n", extensions)
	}

	content, err := collectMediaContent(path)
	if err != nil {
		log.Fatalf("Error during getting files: %s\n", err)
	}
	fmt.Printf("len content: %d\n", len(content))
    tl := buildPlayList(content)
    err = writePlayList(tl)
    if err != nil {
        log.Fatalf("Error during writing playlist file: %s\n", err)
    }
}
