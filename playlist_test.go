package main

import (
	"testing"
)

func TestPlaylist_allocate(t *testing.T) {
	durB := DurationBucket{}
	durations := []float64{3.46, 7.86, 5.67, 8.52, 15, 20, 36.12, 62.45, 70}
	for _, duration := range durations {
		durB.allocate(duration)
	}
	tests := []struct {
		name     string
		got      int
		expected int
	}{
		{"Dur0_5", durB.Dur0_5, 1},
		{"Dur5_10", durB.Dur5_10, 3},
		{"Dur10_30", durB.Dur10_30, 2},
		{"Dur30_60", durB.Dur30_60, 1},
		{"DurOver60", durB.DurOver60, 2},
	}
	for _, tt := range tests {
		if tt.got != tt.expected {
			t.Errorf("Failed %s: got %d, expected %d\n", tt.name, tt.got, tt.expected)
		}
	}
}

func TestPlaylistGetDirRootHasSeparator(t *testing.T) {
	root := "/home/user/Music/"
	expected := "Genre/Artist/Album"
	item := MediaItem{
		AbsPath: "/home/user/Music/Genre/Artist/Album/Track01.mp4",
	}
	got := item.getDir(root)
	if got != expected {
		t.Errorf("Got %s, expected %s\n", got, expected)
	}
}

func TestPlaylistGetDirRootNoSeparator(t *testing.T) {
	root := "/home/user/Music"
	expected := "Genre/Artist/Album"
	item := MediaItem{
		AbsPath: "/home/user/Music/Genre/Artist/Album/Track01.mp4",
	}
	got := item.getDir(root)
	if got != expected {
		t.Errorf("Got %s, expected %s\n", got, expected)
	}
}

func TestPlaylistGetDirRoot(t *testing.T) {
	root := "/home/user/Music/Genre/Artist/Album"
	expected := "Album"
	item := MediaItem{
		AbsPath: "/home/user/Music/Genre/Artist/Album/Track01.mp4",
	}
	got := item.getDir(root)
	if got != expected {
		t.Errorf("Got %s, expected %s\n", got, expected)
	}
}
