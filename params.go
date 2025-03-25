package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	playListExtension = ".xspf"
)

func validateRatio(ratio int) error {
	if ratio < 0 || ratio > 100 {
		return fmt.Errorf("Ratio should be between 0 and 100, got %d\n", ratio)
	}
	return nil
}

func parseParam(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

type Params struct {
	extFlag     bool
	playFlag    bool
	minDuration int
	maxDuration int
	stabilizer  int
	ratio       int
	fileName    string
	includeF    []string
	skipF       []string
	options     Options
	fdate       time.Time
	tdate       time.Time
	optFile     string
	marqueeOpts MarqueeOpts
}

type Options struct {
	Audio     bool
	StartTime uint16
	StopTime  uint16
	Text      string
}

func (o *Options) StringifyAudio() string {
	if !o.Audio {
		return "no-audio"
	}
	return ""
}

func (o *Options) StringifyStartTime() string {
	return "start-time=" + strconv.Itoa(int(o.StartTime))
}

func (o *Options) StringifyEndTime() string {
	return "stop-time=" + strconv.Itoa(int(o.StopTime))
}

func (o *Options) SetSeconds(field string, opt string) error {
	parts := strings.Split(opt, "=")
	seconds, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return fmt.Errorf("Int conversion failed for: %s\n", opt)
	}
	structValue := reflect.ValueOf(o).Elem()
	fieldValue := structValue.FieldByName(field)
	if !fieldValue.IsValid() {
		return fmt.Errorf("field error: %s", field)
	}
	if !fieldValue.CanSet() {
		return fmt.Errorf("cannot set field: %s", field)
	}
	val := reflect.ValueOf(uint16(seconds))
	fieldValue.Set(val)
	return nil
}

func (o *Options) SetText(opt string) error {
	parts := strings.SplitN(opt, "=", 2)
	if len(parts) < 2 || parts[1] == "" {
		return fmt.Errorf("Text option missing value: %s\n", opt)
	}
	text := strings.ReplaceAll(parts[1], `\n`, "\n")
	o.Text = text
	return nil
}

func (p *Params) setFileName(fn string) error {
	if fn == "" {
		p.fileName = "pl-test" + playListExtension
	} else {
		ext := filepath.Ext(fn)
		if ext != "" {
			return fmt.Errorf("File name should not have extension defined")
		}
		p.fileName = fn + playListExtension
	}
	return nil
}

func (p *Params) setDateParams(fdate, tdate string) error {
	fd, err := time.Parse("20060102", fdate)
	if err != nil {
		return fmt.Errorf("Invalid date format - needs YYYYMMDD, got %s for fdate\n", fdate)
	}
	td, err := time.Parse("20060102", tdate)
	if err != nil {
		return fmt.Errorf("Invalid date format - needs YYYYMMDD, got %s for tdate\n", tdate)
	}
	if fd.After(td) {
		return fmt.Errorf("fdate should be before tdate: %v - %v\n", fd, td)
	}
	p.fdate = fd
	p.tdate = td
	return nil
}

func (p *Params) setFolderParams(includeF, skipF string) error {
	if includeF != "" && skipF != "" {
		return fmt.Errorf("Include and skip folders are mutually exclusive")
	}
	p.includeF = parseParam(includeF)
	p.skipF = parseParam(skipF)
	return nil
}

func (p *Params) setOptions(options string) error {
	opts := parseParam(options)
	p.options = Options{Audio: true}
	for _, opt := range opts {
		switch {
		case opt == "no-audio":
			p.options.Audio = false
		case strings.HasPrefix(opt, "start-time"):
			err := p.options.SetSeconds("StartTime", opt)
			if err != nil {
				return fmt.Errorf("Error setting start-time: %s\n", opt)
			}
		case strings.HasPrefix(opt, "stop-time"):
			err := p.options.SetSeconds("StopTime", opt)
			if err != nil {
				return fmt.Errorf("Int conversion failed for end-time: %s\n", opt)
			}
		case strings.HasPrefix(opt, "text"):
			err := p.options.SetText(opt)
			if err != nil {
				return fmt.Errorf("Error setting text: %s\n", opt)
			}
		default:
			return fmt.Errorf("Unrecognized option: %s\n", opt)
		}
	}
	return nil
}

func (p *Params) parseOptFile(fsys fs.FS, fn string) error {
	if fn == "" {
		return nil
	}
	file, err := fsys.Open(fn)
	if err != nil {
		return fmt.Errorf("File cannot be opened: %s", fn)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("File cannot be read: %s", fn)
	}
	var m MarqueeOpts
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("Error unmarshalling options file: %s", err)
	}
	found := m.Marquee.validateColor()
	if !found {
		return fmt.Errorf("Color %s is Unrecognized\n", m.Marquee.Color)
	}
	found = m.Marquee.validatePosition()
	if !found {
		return fmt.Errorf("Position %s is Unrecognized\n", m.Marquee.Position)
	}
	p.marqueeOpts = m
	return nil
}

func getParams() (*Params, error) {
	p := &Params{}
	flag.BoolVar(&p.extFlag, "ext", false, "If specified, collects unique file extensions")
	flag.BoolVar(&p.playFlag, "play", false, "If specified, playlist will be played")
	flag.IntVar(&p.minDuration, "mindur", 0, "Minimum duration of media files to collect (in seconds)")
	flag.IntVar(&p.maxDuration, "maxdur", math.MaxInt32, "Maximum duration of media files to collect (in seconds)")
	flag.IntVar(&p.stabilizer, "stabilizer", math.MaxInt32, "Specifies the interval at which elements are fixed in place during shuffling (they still could be swapped)")
	flag.IntVar(&p.ratio, "ratio", 100, "Specifies the ratio of files to be included in the playlist (e.g. 80 means roughly 80%)")
	fileName := flag.String("fn", "", "Specifies the file name of the playlist")
	fdate := flag.String("fdate", "20000101", "Files created after fdate will be considered")
	tdate := flag.String("tdate", "20300101", "Files created before tdate will be considered")
	includeF := flag.String("include", "", "Folders to consider")
	skipF := flag.String("skip", "", "Folders to skip")
	options := flag.String("options", "", "Options to use:start-time, stop-time, no-audio")
	optFile := flag.String("opt_file", "", "File to set options")
	flag.Parse()
	err := validateRatio(p.ratio)
	if err != nil {
		return nil, err
	}
	err = p.setFileName(*fileName)
	if err != nil {
		return nil, err
	}
	err = p.setDateParams(*fdate, *tdate)
	if err != nil {
		return nil, err
	}
	err = p.setFolderParams(*includeF, *skipF)
	if err != nil {
		return nil, err
	}
	err = p.setOptions(*options)
	if err != nil {
		return nil, err
	}
	fsys := os.DirFS(".")
	err = p.parseOptFile(fsys, *optFile)
	if err != nil {
		return nil, err
	}
	return p, nil
}
