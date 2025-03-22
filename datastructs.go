package main

import "encoding/xml"

var mediaExtensions = []string{".mp4", ".mkv", ".avi", ".flv", ".mpeg"}

const (
	ExtensionApplication = "http://www.videolan.org/vlc/playlist/0"
	Xmlns                = "http://xspf.org/ns/0/"
	XmlnsVlc             = "http://www.videolan.org/vlc/playlist/ns/0/"
)

type Extension struct {
	XMLName     xml.Name `xml:"extension"`
	Application string   `xml:"application,attr"`
	Id          int      `xml:"vlc:id"`
	Options     []string `xml:"vlc:option,omitempty"`
}

type Track struct {
	XMLName  xml.Name  `xml:"track"`
	Location string    `xml:"location"`
	Title    string    `xml:"title"`
	Duration float64   `xml:"duration"`
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

type MarqueeOpts struct {
	Marquee Marquee `json:"marquee"`
}

type Marquee struct {
	Text     string `json:"text,omitempty"`
	Opacity  int    `json:"opacity,omitempty"`
	Position string `json:"position,omitempty"`
}
