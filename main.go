package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"math"
	"time"
)

func randomizePlaylist(playlist []MediaItem) {
    rand.Shuffle(len(playlist), func(i, j int) {
        playlist[i], playlist[j] = playlist[j], playlist[i]
    })
    for i, m := range playlist {
        fmt.Printf("%d -- media id: %d -- media name: %s\n", i, m.Id, m.Name)
    }
}

func main() {
	defer TimeTrack(time.Now(), "main")
	extFlag := flag.Bool("ext", false, "If specified, collects unique file extensions")
	minDuration := flag.Int("mindur", 0, "Minimum duration of media files to collect (in seconds)")
	maxDuration := flag.Int("maxdur", math.MaxInt32, "Maximum duration of media files to collect (in seconds)")
	flag.Parse()

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

	content, err := collectMediaContent(path, *minDuration, *maxDuration)
	if err != nil {
		log.Fatalf("Error during getting files: %s\n", err)
	}
	fmt.Printf("len content: %d\n", len(content))
    randomizePlaylist(content)
	tl := buildPlayList(content)
	err = writePlayList(tl)
	if err != nil {
		log.Fatalf("Error during writing playlist file: %s\n", err)
	}
}
