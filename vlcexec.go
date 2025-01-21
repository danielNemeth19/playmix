package main

import (
	"fmt"
	"image/color"
	"log"

	vlc "github.com/adrg/libvlc-go/v3"
)

func playMixExec(fileName string) {
	// fmt.Printf("playing now: %s\n", fileName)
	// cmd := exec.Command("vlc", fileName, "--sub-filter=marq{marquee=test,color=0x00FF00,position=8}")

	// if err := cmd.Run(); err != nil {
		// log.Fatal(err)
	// }

	fmt.Println(fileName)
	if err := vlc.Init("--fullscreen"); err != nil {
		log.Fatal(err)
	}
	defer vlc.Release()

	// Create a new player.
	player, err := vlc.NewPlayer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		player.Stop()
		player.Release()
	}()

	media, err := player.LoadMediaFromPath("test.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer media.Release()

	marquee := player.Marquee()
	marquee.Enable(true)
	marquee.SetText("TESTING TESTING\nNew Line\nOne more\nLast")
	color := color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}
	marquee.SetColor(color)
	marquee.SetOpacity(100)
	marquee.SetPosition(vlc.PositionBottomRight)
	marquee.SetText("TESTING TESTING\nNew Line\nOne more\nLast")

	// Retrieve player event manager.
	manager, err := player.EventManager()
	if err != nil {
		log.Fatal(err)
	}

	// Register the media end reached event with the event manager.
	quit := make(chan struct{})
	eventCallback := func(event vlc.Event, userData interface{}) {
		close(quit)
	}

	eventID, err := manager.Attach(vlc.MediaPlayerEndReached, eventCallback, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer manager.Detach(eventID)

	// Start playing the media.
	if err = player.Play(); err != nil {
		log.Fatal(err)
	}

	<-quit // }
}
