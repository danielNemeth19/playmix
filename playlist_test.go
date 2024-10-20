package main

import (
	"fmt"
	"playmix/internal/assert"
	"playmix/internal/mocks"
	"testing"
	"testing/fstest"
	"time"
)

func TestDateFilterIn(t *testing.T) {
	params := Params{
		fdate: time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC),
		tdate: time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC),
	}
	inDate := time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC)
	fd := mocks.CreateFakeDirEntry("myfile.mp4", false, inDate)

	got := dateFilter(fd, params)
	assert.Equal(t, "Should be selected", got, true)
}

func TestDateFilterBefore(t *testing.T) {
	params := Params{
		fdate: time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC),
		tdate: time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC),
	}
	beforeDate := time.Date(2022, 6, 16, 0, 0, 0, 0, time.UTC)
	fdBefore := mocks.CreateFakeDirEntry("1.mp4", false, beforeDate)
	got := dateFilter(fdBefore, params)
	assert.Equal(t, "Should not be selected as file is before", got, false)
}
func TestDateFilterAfter(t *testing.T) {
	params := Params{
		fdate: time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC),
		tdate: time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC),
	}
	afterDate := time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC)
	fdAfter := mocks.CreateFakeDirEntry("2.mp4", false, afterDate)
	got := dateFilter(fdAfter, params)
	assert.Equal(t, "Should not be selected as file is after", got, false)
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

// type FSFileOpener struct {
// Files map[string]fstest.MapFile
// }

// func (fs FSFileOpener) Open(fn string) (*os.File, error) {
// mapFile, ok := fs.Files[fn]
// if !ok {
// return nil, fmt.Errorf("File %s not found\n", fn)
// }
// f, err := os.CreateTemp("", "temf")
// if err != nil {
// return nil, fmt.Errorf("Error creating temp file: %v\n", err)
// }
// _, err = f.Write(mapFile.Data)
// if err != nil {
// return nil, fmt.Errorf("Write failed: %v\n", err)
// }
// return f, nil
// }

func TestGetDuration(t *testing.T) {
	fn := "sdfsdftrack.mp4"
	f := fstest.MapFS{
		fn: {
			Data:    []byte("this is my data"),
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	getDurationFS(f, fn)

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

func TestCollectMediaContent(t *testing.T) {
	modTime := time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC)
	fdate := time.Date(2022, 3, 26, 0, 0, 0, 0, time.UTC)
	tdate := time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC)
	fsys := fstest.MapFS{
		"home/Music/track1.mp4": {
			Mode:    0755,
			ModTime: modTime,
		},
	}
	params := Params{fdate: fdate, tdate: tdate, ratio: 100}
	got, _, err := collectMediaContent("home/Music", fsys, params)
	assert.ErrorRaised(t, "Should not raise", err, false)
	fmt.Println(got)
	want := []MediaItem{}
	assert.EqualSlice(t, "Should collect unique extensions", got, want)
}
