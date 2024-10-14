package main

import (
	"playmix/internal/assert"
	"testing"
	"testing/fstest"
	"time"
)

func TestGetPath(t *testing.T) {
	t.Setenv("MEDIA_SOURCE", "/home/user/Music/")
	got, err := getPath()
	expected := "/home/user/Music/"
	assert.Equal(t, "Should parse root path", got, expected)
	assert.ErrorRaised(t, "No error", err, false)
}

func TestGetPathNormalized(t *testing.T) {
	t.Setenv("MEDIA_SOURCE", "/home/user/Music")
	got, err := getPath()
	expected := "/home/user/Music/"
	assert.Equal(t, "Should parse root path normalized", got, expected)
	assert.ErrorRaised(t, "No error", err, false)
}

func TestGetPathRaisesError(t *testing.T) {
	t.Setenv("MEDIA_SOURCE", "")
	_, err := getPath()
	assert.ErrorRaised(t, "Should raise error", err, true)
}

func TestGetPathPartsWithFile(t *testing.T) {
	p := "/home/user/Music/Genre/Artist/Album/Track01.mp4"
	expected := []string{"home", "user", "Music", "Genre", "Artist", "Album"}
	got := getPathParts(p)
	assert.EqualSlice(t, "Should parse path parts - file component", got, expected)
	assert.Equal(t, "Should have parts", len(got), 6)
}

func TestGetPathPartsForRoot(t *testing.T) {
	p := "/home/user/Music/"
	expected := []string{"home", "user", "Music"}
	got := getPathParts(p)
	assert.EqualSlice(t, "Should parse path parts - file component", got, expected)
	assert.Equal(t, "Should have parts", len(got), 3)
}

func TestIsMediaFile(t *testing.T) {
	tests := []struct {
		ext      string
		expected bool
		name     string
	}{
		{
			ext:      ".mp4",
			expected: true,
			name:     "Should be selected: .mp4",
		},
		{
			ext:      ".nfo",
			expected: false,
			name:     "Should not be selected: .nfo",
		},
		{
			ext:      ".mkv",
			expected: true,
			name:     "Should be selected: .mkv",
		},
	}
	for _, tt := range tests {
		verdict := isMediaFile(tt.ext)
		assert.Equal(t, tt.name, verdict, tt.expected)
	}
}

func TestCollectExtensions(t *testing.T) {
	modTime := time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC)
	fsys := fstest.MapFS{
		"home/Music/Album/Artist/best_track.wav": {
			Mode:    0755,
			ModTime: modTime,
		},
		"home/Music/track.mp3": {
			Mode:    0755,
			ModTime: modTime,
		},
		"home/Music/other_track.mp4": {
			Mode:    0755,
			ModTime: modTime,
		},
	}
	got, err := collectExtensions(fsys)
	assert.ErrorRaised(t, "Should not raise", err, false)
	want := []string{".wav", ".mp4", ".mp3"}
	assert.EqualSlice(t, "Should collect unique extensions", got, want)
}
