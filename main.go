package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

const (
	ExtensionApplication = "http://www.videolan.org/vlc/playlist/0"
	Xmlns                = "http://xspf.org/ns/0/"
	XmlnsVlc             = "http://www.videolan.org/vlc/playlist/ns/0/"
)

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

func getPath() string {
	path := os.Getenv("MEDIA_SOURCE")
	if path == "" {
		log.Fatal("MEDIA_SOURCE is not set")
	}
	return path
}

func dumpConsole(s any) {
	out, _ := xml.MarshalIndent(s, " ", "  ")
	fmt.Println(xml.Header + string(out))
}

func writePlayList(s any) {
    outFile, err := os.Create("temp.xspf")
	if err != nil {
        log.Fatalf("Error creating file: %s\n", err)
	}
    outFile.WriteString(xml.Header)
    encoder := xml.NewEncoder(outFile)
    encoder.Indent("", "\t")
    err = encoder.Encode(&s)
	if err != nil {
		log.Fatalf("Error in encoding xml: %s\n", err)
	}
}

func main() {
	path := getPath()
	trackList := buildTrackList(path)

	pl := &PlayList{Xmlns: Xmlns, XmlnsVlc: XmlnsVlc, Version: "1", Title: "Test playlist", Tl: *trackList}
	writePlayList(pl)
}
