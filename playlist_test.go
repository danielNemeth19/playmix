package main

import (
	"io/fs"
	"playmix/internal/assert"
	"testing"
	"time"
)

type MockFileInfo struct {
	fileName string
	size     int64
	mode     fs.FileMode
	modTime  time.Time
	sys      any
}

func (m MockFileInfo) Name() string {
	return m.fileName
}

func (m MockFileInfo) Size() int64 {
	return m.size
}

func (m MockFileInfo) Mode() fs.FileMode {
	return m.mode
}

func (m MockFileInfo) ModTime() time.Time {
	return m.modTime
}

func (m MockFileInfo) IsDir() bool {
	return false
}

func (m MockFileInfo) Sys() any {
	return m.sys
}

type MockDirEntry struct {
	name  string
	mInfo fs.FileInfo
}

func (m MockDirEntry) Name() string {
	return m.name
}

func (m MockDirEntry) IsDir() bool {
	return false
}

func (m MockDirEntry) Type() fs.FileMode {
	return 2
}

func (m MockDirEntry) Info() (fs.FileInfo, error) {
	return m.mInfo, nil
}

func TestDateFilter(t *testing.T) {
	params := Params{
		fdate: time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC),
		tdate: time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC),
	}
	myM := MockFileInfo{
		fileName: "myfile.mp4",
		modTime:  time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	m := MockDirEntry{
		name:  myM.Name(),
		mInfo: myM,
	}
	got := dateFilter(m, params)
	assert.Equal(t, "Should be selected", got, true)
}

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

func TestIsIncludedIfEmptyFilter(t *testing.T) {
	root := []string{"home", "user", "Music"}
	path := "/home/user/Music/Genre/Artist/Album/Track01.mp4"
	inc := []string{}
	got := isIncluded(root, path, inc)
	assert.Equal(t, "Should be included", got, true)
}

func TestIsIncluded(t *testing.T) {
	root := []string{"home", "user", "Music"}
	path := "/home/user/Music/Genre/Artist/Album/Track01.mp4"
	inc := []string{"Album"}
	got := isIncluded(root, path, inc)
	assert.Equal(t, "Should be included", got, true)
}

func TestIsIncludedFalse(t *testing.T) {
	root := []string{"home", "user", "Music"}
	path := "/home/user/Music/Genre/Artist/Album/Track01.mp4"
	inc := []string{"OtherArtist"}
	got := isIncluded(root, path, inc)
	assert.Equal(t, "Should not be included", got, false)
}

func TestIsIncludedFalseIfFolderNotWithinRoot(t *testing.T) {
	root := []string{"home", "user", "Music"}
	path := "/home/user/Music/Genre/Artist/Album/Track01.mp4"
	inc := []string{"home", "user"}
	got := isIncluded(root, path, inc)
	assert.Equal(t, "Should not be included", got, false)
}

func TestPlaylistGetDirRoot(t *testing.T) {
	root := []string{"home", "user", "Music"}
	expected := "Music"
	item := MediaItem{
		AbsPath: "/home/user/Music/Track01.mp4",
	}
	got := item.getRelativeDir(root)
	assert.Equal(t, "Should get relative dir for file in root", got, expected)
}

func TestPlaylistGetDirSubFolders(t *testing.T) {
	root := []string{"home", "user", "Music"}
	expected := "Genre/Artist/Album"
	item := MediaItem{
		AbsPath: "/home/user/Music/Genre/Artist/Album/Track01.mp4",
	}
	got := item.getRelativeDir(root)
	assert.Equal(t, "Should get relative dir for file in subfolder", got, expected)
}
