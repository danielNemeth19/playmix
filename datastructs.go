package main

import (
	"encoding/xml"
	vlc "github.com/adrg/libvlc-go/v3"
)

var mediaExtensions = []string{".mp4", ".mkv", ".avi", ".flv", ".mpeg"}

const (
	ExtensionApplication = "http://www.videolan.org/vlc/playlist/0"
	Xmlns                = "http://xspf.org/ns/0/"
	XmlnsVlc             = "http://www.videolan.org/vlc/playlist/ns/0/"
)

var textPositionMap = map[string]vlc.Position{
	"disable":     vlc.PositionDisable,
	"center":      vlc.PositionCenter,
	"left":        vlc.PositionLeft,
	"right":       vlc.PositionRight,
	"top":         vlc.PositionTop,
	"topleft":     vlc.PositionTopLeft,
	"topright":    vlc.PositionTopRight,
	"bottom":      vlc.PositionBottom,
	"bottomleft":  vlc.PositionBottomLeft,
	"bottomright": vlc.PositionBottomRight,
}

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

func (m Marquee) validatePosition() bool {
	_, found := textPositionMap[m.Position]
	return found
}

func (m Marquee) remapPosition() vlc.Position {
	position, found := textPositionMap[m.Position]
	if !found {
		return vlc.PositionDisable
	}
	return position
}
