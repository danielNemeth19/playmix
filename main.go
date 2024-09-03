package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"time"
)

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

func main() {
	defer TimeTrack(time.Now(), "main")
	extFlag := flag.Bool("ext", false, "If specified, collects unique file extensions")
	minDuration := flag.Int("mind", 0, "If specified, collects files with duration more than specified")
	maxDuration := flag.Int("maxd", 60*60, "If specified, collects files with duration less than specified")
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

	content, err := collectMediaContent(path, *minDuration, *maxDuration)
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
