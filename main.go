package main

import (
	"encoding/xml"
	"fmt"
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
	Id          string   `xml:"vlc:id"`
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

func ListFiles() {
	path := "/path/to/files"
	files, _ := os.ReadDir(path)
	for i, v := range files {
		info, _ := v.Info()
		fmt.Printf("key: %d -- name: %s\n", i, info.Name())
	}
}

func main() {
	ext := &Extension{Application: ExtensionApplication, Id: "01"}
	extOut, err := xml.MarshalIndent(ext, "", "  ")
	if err != nil {
		fmt.Printf("fix this later, now error: %s\n", err)
	}
	fmt.Println(string(extOut))

	track1 := &Track{Location: "/mnt/path/", Ext: *ext}
	out, err := xml.MarshalIndent(track1, " ", "  ")
	if err != nil {
		fmt.Printf("fix this later, now error: %s\n", err)
	}
	fmt.Println(string(out))
	track2 := &Track{Location: "/mnt/path/", Ext: Extension{Application: ExtensionApplication, Id: "02"}}

	tl := &TrackList{}
	tl.Tracks = []*Track{track1, track2}

	tlOut, err := xml.MarshalIndent(tl, " ", "  ")
	if err != nil {
		fmt.Printf("fix this later, now error: %s\n", err)
	}
	fmt.Println(string(tlOut))

	pl := &PlayList{Xmlns: Xmlns, XmlnsVlc: XmlnsVlc, Version: "1", Title: "Test playlist", Tl: *tl}
	plOut, err := xml.MarshalIndent(pl, " ", "  ")
	if err != nil {
		fmt.Printf("fix this later, now error: %s\n", err)
	}
	fmt.Println(xml.Header + string(plOut))
	ListFiles()
}
