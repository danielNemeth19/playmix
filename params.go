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
	fdate       string
	tdate       string
}

func validateParams(includeF, skipF string, ratio int) error {
	if includeF != "" && skipF != "" {
		return fmt.Errorf("Include and skip folders are mutually exclusive")
	}
	if ratio < 0 || ratio > 100 {
		return fmt.Errorf("Ratio should be between 0 and 100, got %d\n", ratio)
	}
	return nil
}

func parse(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func (p Params) parseDateString(dateType string) time.Time {
    dateMap := map[string]string{
        "fdate": p.fdate,
        "tdate": p.tdate,
    }
    toGet, ok := dateMap[dateType]
    if !ok {
        return fmt.Errorf("Let's think about this..this might need to be run during validation")
    }
    fmt.Printf("ok? %v\n", ok)
    date, _ := time.Parse("20060102", toGet)
	return date
}

func getParams() (*Params, error) {
	p := &Params{}
	flag.BoolVar(&p.extFlag, "ext", false, "If specified, collects unique file extensions")
	flag.IntVar(&p.minDuration, "mindur", 0, "Minimum duration of media files to collect (in seconds)")
	flag.IntVar(&p.maxDuration, "maxdur", math.MaxInt32, "Maximum duration of media files to collect (in seconds)")
	flag.IntVar(&p.stabilizer, "stabilizer", math.MaxInt32, "Specifies the interval at which elements are fixed in place during shuffling (they still could be swapped)")
	flag.IntVar(&p.ratio, "ratio", 100, "Specifies the ratio of files to be included in the playlist (e.g. 80 means roughly 80%)")
	flag.StringVar(&p.fdate, "fdate", "20000101", "Files created after fdate will be considered")
	flag.StringVar(&p.tdate, "tdate", "20300101", "Files created before tdate will be considered")
	includeF := flag.String("include", "", "Folders to consider")
	skipF := flag.String("skip", "", "Folders to skip")
	flag.Parse()
	err := validateParams(*includeF, *skipF, p.ratio)
	if err != nil {
		return nil, err
	}
	p.includeF = parse(*includeF)
	p.skipF = parse(*skipF)
	return p, nil
}
