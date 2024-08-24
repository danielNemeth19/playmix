package main

import (
	"encoding/xml"
	"flag"
	"fmt"
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

var videoExtensions = []string{".mp4", ".mkv", ".avi", ".flv", ".mpeg"}

type Extension struct {
	XMLName     xml.Name `xml:"extension"`
	Application string   `xml:"application,attr"`
	Id          int      `xml:"vlc:id"`
}

type Track struct {
	XMLName  xml.Name  `xml:"track"`
	Location string    `xml:"location"`
	Ext      Extension `xml:"extension"`
}

type TrackList struct {
	XMLName xml.Name `xml:"trackList"`
	Tracks  []*Track `xml:"track"`
}

type PlayList struct {
	XMLName  xml.Name  `xml:"playlist"`
	Xmlns    string    `xml:"xmlns,attr"`
	XmlnsVlc string    `xml:"xmlns:vlc,attr"`
	Version  string    `xml:"version,attr"`
	Title    string    `xml:"title"`
	Tl       TrackList `xml:"trackList"`
}

type MediaItem struct {
	Name     string
	Id       int
	Duration float64
}

func isVideoFile(ext string) bool {
	for _, v := range videoExtensions {
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
		if isVideoFile(extension) {
			fmt.Println(i, file.Name())
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
	log.Printf("video file to check: %s\n", p)
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
	fmt.Printf("Duration: %d\n", mp4.Moov.Mvhd.Duration)
	fmt.Printf("Time scale: %d\n", mp4.Moov.Mvhd.Timescale)
	fmt.Printf("Duration in seconds: %f\n", float64(mp4.Moov.Mvhd.Duration)/float64(mp4.Moov.Mvhd.Timescale))

    duration := float64(mp4.Moov.Mvhd.Duration) / float64(mp4.Moov.Mvhd.Timescale)
	return duration, nil
}

func main() {
	path, err := getPath()
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	log.Printf("Path to be used: %s\n", path)

	extFlag := flag.Bool("ext", false, "If specified, collects unique file extensions")
	flag.Parse()

	fmt.Printf("Ext flag: %v\n", *extFlag)
	if *extFlag {
		extensions, err := getExtensions(path)
		if err != nil {
			log.Fatalf("Error during extension collection: %s\n", err)
		}
		fmt.Printf("Extensions: %v\n", *extensions)
	}

	content, err := getFolderContent(path)
	if err != nil {
		log.Fatalf("Error during getting files: %s\n", err)
	}
    for _, item := range content {
        fmt.Printf("%d -- %s -- %f\n", item.Id, item.Name, item.Duration)
    }
}
