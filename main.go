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
	
	log.Printf("Path to be used: %s\n", params.MediaPath)

	fsys := os.DirFS(params.MediaPath)
	if params.extFlag {
		extensions, err := collectExtensions(fsys)
		if err != nil {
			log.Fatalf("Error during extension collection: %s\n", err)
		}
		fmt.Printf("Extensions: %v\n", extensions)
	}
	content, summary, err := collectMediaContent(params.MediaPath, fsys, *params)
	if err != nil {
		log.Fatalf("Error during getting files: %s\n", err)
	}
	randomizePlaylist(content, int(params.RandomizerOptions.Stabilizer))
	playList := buildPlayList(content, params.PlayOptions)

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
		playMixList(params.fileName, params.MarqueeOptions)
	}
}
