package main

import (
	"log"

	vlc "github.com/adrg/libvlc-go/v3"
)

func playMix(fileName string) {
	// Initialize libVLC. Additional command line arguments can be passed in
	// to libVLC by specifying them in the Init function.
	if err := vlc.Init("--fullscreen"); err != nil {
	// --sub-filter="marq{marquee=akarmi ami érdekes}"
    // if err := vlc.Init(":sub-filter=marq", ":marq-marquee=Hello, World!"); err != nil {
		log.Fatal(err)
	}
	defer vlc.Release()

	// Create a new list player.
	player, err := vlc.NewListPlayer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		player.Stop()
		player.Release()
	}()

	list, err := vlc.NewMediaList()
	if err != nil {
		log.Fatal(err)
	}
	defer list.Release()

	// err = list.AddMediaFromPath(fileName)
	// if err != nil {
		// log.Fatal(err)
	// }
	
	media, err := vlc.NewMediaFromPath(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer media.Release()


	media.AddOptions("--sub-filter=\"marq{marquee=akarmi ami érdekes}\"")

    if err = list.AddMedia(media); err != nil {
        log.Fatal(err)
    }

    if err = player.SetMediaList(list); err != nil {
		log.Fatal(err)
	}

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

	<-quit
}
