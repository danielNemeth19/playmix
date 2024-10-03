package main

import (
	"playmix/internal/assert"
	"testing"
)

func TestPlaylist_allocate(t *testing.T) {
	durB := DurationBucket{}
	durations := []float64{3.46, 7.86, 8.52, 15, 20, 36.12, 62.45, 70, 190, 241}
	for _, duration := range durations {
		durB.allocate(duration)
	}
	tests := []struct {
		name     string
		got      int
		expected int
	}{
		{"Dur0_5", durB.Dur0_5, 1},
		{"Dur5_10", durB.Dur5_10, 2},
		{"Dur10_30", durB.Dur10_30, 2},
		{"Dur30_60", durB.Dur30_60, 1},
		{"Dur60_180", durB.Dur60_180, 2},
		{"Dur180_240", durB.Dur180_240, 1},
		{"DurOver240", durB.DurOver240, 1},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.name, tt.got, tt.expected)
	}
}

func TestPlaylistToSkip(t *testing.T) {
	skipF := []string{"c", "e"}
	verdict := toSkip("a", skipF)
	assert.Equal(t, "toSkip to be false", verdict, false)
	verdict = toSkip("c", skipF)
	assert.Equal(t, "toSkip to be true", verdict, true)
}

func TestPlaylistGetDirRootHasSeparator(t *testing.T) {
	root := "/home/user/Music/"
	expected := "Genre/Artist/Album"
	item := MediaItem{
		AbsPath: "/home/user/Music/Genre/Artist/Album/Track01.mp4",
	}
	got := item.getDir(root)
	assert.Equal(t, "root has separator", got, expected)
}

func TestPlaylistGetDirRootNoSeparator(t *testing.T) {
	root := "/home/user/Music"
	expected := "Genre/Artist/Album"
	item := MediaItem{
		AbsPath: "/home/user/Music/Genre/Artist/Album/Track01.mp4",
	}
	got := item.getDir(root)
	assert.Equal(t, "root no separator", got, expected)
}

func TestPlaylistGetDirRoot(t *testing.T) {
	root := "/home/user/Music/Genre/Artist/Album"
	expected := "Album"
	item := MediaItem{
		AbsPath: "/home/user/Music/Genre/Artist/Album/Track01.mp4",
	}
	got := item.getDir(root)
	assert.Equal(t, "File in root", got, expected)
}
