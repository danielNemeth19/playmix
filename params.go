package main

import (
	"flag"
	"math"
	"strings"
)

type Params struct {
	extFlag     bool
	minDuration int
	maxDuration int
	stabilizer  int
	ratio       int
	includeF    []string
	skipF       []string
}

func parse(s string) []string {
    if s == "" {
        return []string{}
    } 
    return strings.Split(s, ",")
}

func getParams() *Params {
	p := &Params{}
	flag.BoolVar(&p.extFlag, "ext", false, "If specified, collects unique file extensions")
	flag.IntVar(&p.minDuration, "mindur", 0, "Minimum duration of media files to collect (in seconds)")
	flag.IntVar(&p.maxDuration, "maxdur", math.MaxInt32, "Maximum duration of media files to collect (in seconds)")
	flag.IntVar(&p.stabilizer, "stabilizer", math.MaxInt32, "Specifies the interval at which elements are fixed in place during shuffling (they still could be swapped)")
	flag.IntVar(&p.ratio, "ratio", 100, "Specifies the ratio of files to be included in the playlist (e.g. 80 means roughly 80%)")
	includeF := flag.String("include", "", "Folders to consider")
	skipF := flag.String("skip", "", "Folders to skip")
	flag.Parse()
	p.includeF = parse(*includeF)
	p.skipF = parse(*skipF)
	return p
}
