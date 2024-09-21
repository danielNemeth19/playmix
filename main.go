package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

func randomizePlaylist(playlist []MediaItem, stabilizer int) {
	if stabilizer > len(playlist) {
		rand.Shuffle(len(playlist), func(i, j int) {
			playlist[i], playlist[j] = playlist[j], playlist[i]
		})
	} else {
		rand.Shuffle(len(playlist), func(i, j int) {
			if i%stabilizer != 0 {
				playlist[i], playlist[j] = playlist[j], playlist[i]
			}
		})
	}
}

func main() {
	defer TimeTrack(time.Now(), "main")
	// TODO: filename needs to be generated, or provided
	extFlag := flag.Bool("ext", false, "If specified, collects unique file extensions")
	minDuration := flag.Int("mindur", 0, "Minimum duration of media files to collect (in seconds)")
	maxDuration := flag.Int("maxdur", math.MaxInt32, "Maximum duration of media files to collect (in seconds)")
	stabilizer := flag.Int("stabilizer", math.MaxInt32, "Specifies the interval at which elements are fixed in place during shuffling (they still could be swapped)")
	ratio := flag.Int("ratio", 100, "Specifies the ratio of files to be included in the playlist (e.g. 80 means roughly 80%)")
	folders := flag.String("folders", "", "Folders to consider")
	flag.Parse()

	fmt.Printf("Folder needed: %s\n", *folders)

	path, err := getPath()
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	log.Printf("Path to be used: %s\n", path)

	if *extFlag {
		extensions, err := collectExtensions(path)
		if err != nil {
			log.Fatalf("Error during extension collection: %s\n", err)
		}
		fmt.Printf("Extensions: %v\n", extensions)
	}

	content, summary, err := collectMediaContent(path, *minDuration, *maxDuration, *ratio)
	if err != nil {
		log.Fatalf("Error during getting files: %s\n", err)
	}
	randomizePlaylist(content, *stabilizer)
	tl := buildPlayList(content)
	err = writePlayList(tl)
	if err != nil {
		log.Fatalf("Error during writing playlist file: %s\n", err)
	}
	summary.getData()
}
