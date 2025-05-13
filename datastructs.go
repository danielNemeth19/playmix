package main

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"log"
	"path/filepath"
	"strconv"
	"strings"

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

var colorMap = map[string]color.RGBA{
	"blue":   {0, 0, 255, 255},
	"yellow": {255, 255, 0, 255},
	"red":    {255, 0, 0, 255},
	"black":  {0, 0, 0, 255},
	"cyan":   {0, 255, 255, 255},
	"white":  {255, 255, 255, 255},
	"green":  {0, 255, 0, 255},
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

type FileOptions struct {
	MediaPath         string            `json:"media_path"`
	FileName          string            `json:"file_name"`
	Marquee           Marquee           `json:"marquee"`
	PlayOptions       PlayOptions       `json:"play_options"`
	RandomizerOptions RandomizerOptions `json:"randomizer_options"`
}

func (f *FileOptions) validatePath() error {
	if f.MediaPath == "" {
		return fmt.Errorf("media_path needs to be set in options file")
	}
	if !strings.HasSuffix(f.MediaPath, string(filepath.Separator)) {
		f.MediaPath += string(filepath.Separator)
		log.Printf("Root path got normalized by adding path separator (%s)\n", string(filepath.Separator))
	}
	return nil
}

type Marquee struct {
	Text     string `json:"text,omitempty"`
	Color    string `json:"color,omitempty"`
	Opacity  int    `json:"opacity,omitempty"`
	Position string `json:"position,omitempty"`
}

func (m Marquee) validateColor() error {
	_, found := colorMap[m.Color]
	if !found && m.Color != "" {
		return fmt.Errorf("%s color not found in color map\n", m.Color)
	}
	return nil
}

func (m Marquee) validatePosition() error {
	_, found := textPositionMap[m.Position]
	if !found && m.Position != "" {
		return fmt.Errorf("%s position not found in position map\n", m.Position)
	}
	return nil
}

func (m Marquee) remapColor() color.RGBA {
	color, found := colorMap[m.Color]
	if !found {
		return colorMap["red"]
	}
	return color
}

func (m Marquee) remapPosition() vlc.Position {
	position, found := textPositionMap[m.Position]
	if !found {
		return vlc.PositionDisable
	}
	return position
}

type PlayOptions struct {
	Audio     bool   `json:"audio,omitempty"`
	StartTime uint16 `json:"start_time,omitempty"`
	StopTime  uint16 `json:"stop_time,omitempty"`
}

func (p PlayOptions) validateTimes() error {
	if p.StartTime != 0 && p.StopTime != 0 &&
		(p.StartTime >= p.StopTime) {
		return fmt.Errorf("Stop time (%d) should be greater than start time (%d)\n", p.StopTime, p.StartTime)
	}
	return nil
}

func (p PlayOptions) StringifyAudio() string {
	if !p.Audio {
		return "no-audio"
	}
	return ""
}

func (p PlayOptions) StringifyStartTime() string {
	return "start-time=" + strconv.Itoa(int(p.StartTime))
}

func (p PlayOptions) StringifyStopTime() string {
	return "stop-time=" + strconv.Itoa(int(p.StopTime))
}

type RandomizerOptions struct {
	Ratio      uint8  `json:"ratio,omitempty"`
	Stabilizer uint32 `json:"stabilizer,omitempty"`
}

func (r RandomizerOptions) validateRatio() error {
	if r.Ratio < 0 || r.Ratio > 100 {
		return fmt.Errorf("Ratio should be between 0 and 100, got %d\n", r.Ratio)
	}
	return nil
}

// TODO: think about this - defaulting is confusing
func (r *RandomizerOptions) setDefaultRatio() {
	if r.Ratio == 0 {
		r.Ratio = 100
	}
}
