package main

import (
	"flag"
	"fmt"
	"math"
	"strings"
	"time"
)

type Params struct {
	extFlag     bool
	minDuration int
	maxDuration int
	stabilizer  int
	ratio       int
	includeF    []string
	skipF       []string
	fdate       time.Time
	tdate       time.Time
}

func validateRatio(ratio int) error {
	if ratio < 0 || ratio > 100 {
		return fmt.Errorf("Ratio should be between 0 and 100, got %d\n", ratio)
	}
	return nil
}

func parseFolder(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
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
    p.fdate = fd
    p.tdate = td
	return nil
}

func (p *Params) setFolderParams(includeF, skipF string) error {
	if includeF != "" && skipF != "" {
		return fmt.Errorf("Include and skip folders are mutually exclusive")
	}
	p.includeF = parseFolder(includeF)
	p.skipF = parseFolder(skipF)
    return nil
}

func getParams() (*Params, error) {
	p := &Params{}
	flag.BoolVar(&p.extFlag, "ext", false, "If specified, collects unique file extensions")
	flag.IntVar(&p.minDuration, "mindur", 0, "Minimum duration of media files to collect (in seconds)")
	flag.IntVar(&p.maxDuration, "maxdur", math.MaxInt32, "Maximum duration of media files to collect (in seconds)")
	flag.IntVar(&p.stabilizer, "stabilizer", math.MaxInt32, "Specifies the interval at which elements are fixed in place during shuffling (they still could be swapped)")
	flag.IntVar(&p.ratio, "ratio", 100, "Specifies the ratio of files to be included in the playlist (e.g. 80 means roughly 80%)")
    fdate := flag.String("fdate", "20000101", "Files created after fdate will be considered")
    tdate := flag.String("tdate", "20300101", "Files created before tdate will be considered")
	includeF := flag.String("include", "", "Folders to consider")
	skipF := flag.String("skip", "", "Folders to skip")
	flag.Parse()
	err := validateRatio(p.ratio)
	if err != nil {
		return nil, err
	}
    err = p.setFolderParams(*includeF, *skipF)
    if err != nil {
        return nil, err
    }
    err = p.setDateParams(*fdate, *tdate)
    if err != nil {
        return nil, err
    }
	return p, nil
}
