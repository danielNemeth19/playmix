package main

import (
	"fmt"
	"log"
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
	params := getParams()
	path, err := getPath()
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	log.Printf("Path to be used: %s\n", path)

	if params.extFlag {
		extensions, err := collectExtensions(path)
		if err != nil {
			log.Fatalf("Error during extension collection: %s\n", err)
		}
		fmt.Printf("Extensions: %v\n", extensions)
	}
	content, summary, err := collectMediaContent(path, *params)
	if err != nil {
		log.Fatalf("Error during getting files: %s\n", err)
	}
	randomizePlaylist(content, params.stabilizer)
	tl := buildPlayList(content)
	err = writePlayList(tl)
	if err != nil {
		log.Fatalf("Error during writing playlist file: %s\n", err)
	}
	summary.getData()
}
