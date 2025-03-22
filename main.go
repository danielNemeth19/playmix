package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	defer TimeTrack(time.Now(), "main")
	params, err := getParams()
	if err != nil {
		log.Fatalf("Param validation error: %s\n", err)
	}
	root, err := getPath()
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	log.Printf("Path to be used: %s\n", root)

	fsys := os.DirFS(root)
	if params.extFlag {
		extensions, err := collectExtensions(fsys)
		if err != nil {
			log.Fatalf("Error during extension collection: %s\n", err)
		}
		fmt.Printf("Extensions: %v\n", extensions)
	}
	content, summary, err := collectMediaContent(root, fsys, *params)
	if err != nil {
		log.Fatalf("Error during getting files: %s\n", err)
	}
	randomizePlaylist(content, params.stabilizer)
	playList := buildPlayList(content, params.options)

	outfile, err := createFile(params.fileName)
	if err != nil {
		log.Fatalf("Error during creating file: %s\n", err)
	}
	defer outfile.Close()

	err = writePlayList(playList, outfile)
	if err != nil {
		log.Fatalf("Error during writing playlist file: %s\n", err)
	}
	// TODO: maybe make duration bucket summary optional too
	summary.getData(os.Stdout)
	if params.playFlag {
		playMixList(params.fileName, params.marqueeOpts.Marquee)
	}
}
