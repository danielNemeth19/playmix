package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"time"
)

const (
	playListExtension = ".xspf"
)

type Params struct {
	extFlag           bool
	playFlag          bool
	minDuration       int
	maxDuration       int
	fdate             time.Time
	tdate             time.Time
	optFile           string
	MediaPath         string
	FileName          string
	MarqueeOptions    Marquee
	PlayOptions       PlayOptions
	RandomizerOptions RandomizerOptions
	FilterOptions     FilterOptions
}

func (p *Params) setFileName(fn string) error {
	if fn == "" {
		p.FileName = "pl-test" + playListExtension
	} else {
		ext := filepath.Ext(fn)
		if ext != "" {
			return fmt.Errorf("File name should not have extension defined")
		}
		p.FileName = fn + playListExtension
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

// func (p *Params) setFolderParams(includeF, skipF string) error {
	// if includeF != "" && skipF != "" {
		// return fmt.Errorf("Include and skip folders are mutually exclusive")
	// }
	// p.includeF = parseParam(includeF)
	// p.skipF = parseParam(skipF)
	// return nil
// }

// TODO: Think about refactoring option group validations out
func (p *Params) parseOptFile(fsys fs.FS, fn string) error {
	data, err := readInOptFile(fsys, fn)
	if err != nil {
		return fmt.Errorf("Something wrong: %s", err)
	}
	var opt FileOptions
	err = json.Unmarshal(data, &opt)
	if err != nil {
		return fmt.Errorf("Error unmarshalling options file: %s", err)
	}
	p.MediaPath = opt.MediaPath
	p.MarqueeOptions = opt.Marquee
	p.PlayOptions = opt.PlayOptions
	p.RandomizerOptions = opt.RandomizerOptions
	p.FilterOptions = opt.FilterOptions

	err = opt.validatePath()
	if err != nil {
		return err
	}

	err = p.setFileName(opt.FileName)
	if err != nil {
		return err
	}
	
	err = p.FilterOptions.validateFilterOptions()
	if err != nil {
		return err
	}

	err = p.MarqueeOptions.validateColor()
	if err != nil {
		return err
	}
	err = p.MarqueeOptions.validatePosition()
	if err != nil {
		return err
	}

	err = p.PlayOptions.validateTimes()
	if err != nil {
		return err
	}
	err = p.RandomizerOptions.validateRatio()
	if err != nil {
		return err
	}
	p.RandomizerOptions.setDefaultRatio()
	return nil
}

func getParams() (*Params, error) {
	p := &Params{}
	flag.BoolVar(&p.extFlag, "ext", false, "If specified, collects unique file extensions")
	flag.BoolVar(&p.playFlag, "play", false, "If specified, playlist will be played")
	flag.IntVar(&p.minDuration, "mindur", 0, "Minimum duration of media files to collect (in seconds)")
	flag.IntVar(&p.maxDuration, "maxdur", math.MaxInt32, "Maximum duration of media files to collect (in seconds)")
	fdate := flag.String("fdate", "20000101", "Files created after fdate will be considered")
	tdate := flag.String("tdate", "20300101", "Files created before tdate will be considered")
	// includeF := flag.String("include", "", "Folders to consider")
	// skipF := flag.String("skip", "", "Folders to skip")
	optFile := flag.String("opt_file", "", "File to set options")
	flag.Parse()
	err := p.setDateParams(*fdate, *tdate)
	if err != nil {
		return nil, err
	}
	// err = p.setFolderParams(*includeF, *skipF)
	// if err != nil {
		// return nil, err
	// }
	fsys := os.DirFS(".")
	if *optFile != "" {
		err = p.parseOptFile(fsys, *optFile)
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}
