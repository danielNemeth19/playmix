package main

import (
	"image/color"
	"log"

	vlc "github.com/adrg/libvlc-go/v3"
)

func playMixList(fileName string) {
	if err := vlc.Init("--fullscreen"); err != nil {
		log.Fatal(err)
	}
	defer vlc.Release()

	// Create a new list listPlayer.
	listPlayer, err := vlc.NewListPlayer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		listPlayer.Stop()
		listPlayer.Release()
	}()

	mediaList, err := vlc.NewMediaList()
	if err != nil {
		log.Fatal(err)
	}
	defer mediaList.Release()

	err = mediaList.AddMediaFromPath(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if err = listPlayer.SetMediaList(mediaList); err != nil {
		log.Fatal(err)
	}

	player, err := listPlayer.Player()
	if err != nil {
		log.Fatal(err)
	}

	// retrieve player instance event manager
	playerManager, err := player.EventManager()
	if err != nil {
		log.Fatal(err)
	}

	marqueeSetter := func(event vlc.Event, userData interface{}) {
		player, err := listPlayer.Player()
		if err != nil {
			log.Fatal(err)
		}
		marquee := player.Marquee()
		marquee.Enable(true)
		marquee.SetText("TESTING TESTING\nTESTING TESTING")
		color := color.RGBA{
			R: 255,
			G: 0,
			B: 0,
			A: 255,
		}
		marquee.SetColor(color)
		marquee.SetOpacity(100)
		marquee.SetPosition(vlc.PositionBottomLeft)
	}

	playerEventID, err := playerManager.Attach(vlc.MediaPlayerPlaying, marqueeSetter, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer playerManager.Detach(playerEventID)

	// Retrieve list player event manager.
	manager, err := listPlayer.EventManager()
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
	if err = listPlayer.Play(); err != nil {
		log.Fatal(err)
	}
	<-quit
}
