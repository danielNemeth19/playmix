package main

import (
	"bytes"
	"math"
	"playmix/internal/assert"
	"playmix/internal/mocks"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

func TestGetRealRation(t *testing.T) {
	s := Summarizer{
		totalScanned:  100,
		totalSelected: 10,
	}
	got := s.getRealRatio()
	assert.Equal(t, "Should be 10%", got, 10)
}

func TestGetData(t *testing.T) {
	d := DurationBucket{
		Dur0_5:     1,
		Dur5_10:    2,
		Dur10_30:   3,
		Dur30_60:   4,
		Dur60_180:  5,
		Dur180_240: 6,
		DurOver240: 7,
	}
	s := Summarizer{totalScanned: 100, totalSelected: 100, dBucket: d, totalDuration: 300, ratio: 100}
	var buf bytes.Buffer
	s.getData(&buf)
	res := buf.String()
	splits := strings.Split(res, "\n")

	assert.Equal(t, "Should be equal", splits[0], "Total scanned: 100")
	assert.Equal(t, "Should be equal", splits[1], "Duration distribution:")
	assert.Equal(t, "Should be equal", splits[2], "Bucket <5 seconds: 1")
	assert.Equal(t, "Should be equal", splits[3], "Bucket 5-10 seconds: 2")
	assert.Equal(t, "Should be equal", splits[4], "Bucket 10-30 seconds: 3")
	assert.Equal(t, "Should be equal", splits[5], "Bucket 30-60 seconds: 4")
	assert.Equal(t, "Should be equal", splits[6], "Bucket 60-180 seconds: 5")
	assert.Equal(t, "Should be equal", splits[7], "Bucket 180-240 seconds: 6")
	assert.Equal(t, "Should be equal", splits[8], "Bucket 240< seconds: 7")
	assert.Equal(t, "Should be equal", splits[9], "Total duration is: 300.000000 sec -- (5.000000) minutes")
	assert.Equal(t, "Should be equal", splits[10], "Total selected: 100 -- required ratio: 100 -- got: 100.00%")
}

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

func TestPlaylistGetRelativeDir(t *testing.T) {
	root := []string{"home", "user", "Music"}
	expected := "Music"
	item := MediaItem{
		AbsPath: "/home/user/Music/Track01.mp4",
	}
	item.getRelativeDir(root)
	assert.Equal(t, "Should get relative dir for file in root", item.Dir, expected)
}

func TestPlaylistGetRelativeDirSubfolder(t *testing.T) {
	root := []string{"home", "user", "Music"}
	expected := "Music/Genre/Artist/Album"
	item := MediaItem{
		AbsPath: "/home/user/Music/Genre/Artist/Album/Track01.mp4",
	}
	item.getRelativeDir(root)
	assert.Equal(t, "Should get relative dir for file in subfolder", item.Dir, expected)
}

func TestGetDuration(t *testing.T) {
	data := mocks.CreateData(60)
	fn := "track.mp4"
	f := fstest.MapFS{
		fn: {
			Data:    data,
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	duration, err := getDuration(f, fn)
	assert.Equal(t, "Should be 60 seconds", duration, 60)
	assert.ErrorRaised(t, "Should not raise error", err, false)
}

func TestGetDurationErrorNonMediaFile(t *testing.T) {
	fn := "track.mp4"
	f := fstest.MapFS{
		fn: {
			Data:    []byte("No Moovbox so not mp4 file"),
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	duration, err := getDuration(f, fn)
	assert.Equal(t, "Should be 0", duration, 0)
	assert.ErrorRaised(t, "Should raise error", err, true)
}

func TestGetDurationErrorOpen(t *testing.T) {
	fn := "../"
	f := fstest.MapFS{
		fn: {
			Data:    []byte(""),
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	duration, err := getDuration(f, fn)
	assert.Equal(t, "Should be 0", duration, 0)
	assert.ErrorRaised(t, "Should raise error if file doesn't exists", err, true)
}

func TestGetDurationErrorStat(t *testing.T) {
	fsys := mocks.FakeSys{}
	duration, err := getDuration(fsys, "stat_error_test")
	assert.Equal(t, "Should return duration as 0", duration, 0)
	assert.ErrorRaised(t, "Should return faked error", err, true)
}

func TestCollectMediaContentRaisesWalkError(t *testing.T) {
	fdate := time.Date(2022, 3, 26, 0, 0, 0, 0, time.UTC)
	tdate := time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC)
	f := mocks.FakeSys{}
	params := Params{fdate: fdate, tdate: tdate, ratio: 100}
	_, _, err := collectMediaContent("/home/Music", f, params)
	assert.ErrorRaised(t, "Should raise error", err, true)
}

func TestCollectMediaContentRaisesErrorNonMediaFile(t *testing.T) {
	modTime := time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC)
	fdate := time.Date(2022, 3, 26, 0, 0, 0, 0, time.UTC)
	tdate := time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC)
	fsys := fstest.MapFS{
		"track1.mp4": {
			Mode:    0755,
			ModTime: modTime,
		},
	}
	params := Params{fdate: fdate, tdate: tdate, ratio: 100}
	_, _, err := collectMediaContent("/home/Music", fsys, params)
	assert.ErrorRaised(t, "Should raise error", err, true)
}

func TestCollectMediaContentDuration(t *testing.T) {
	modTime := time.Date(2020, 3, 26, 0, 0, 0, 0, time.UTC)
	fdate := time.Date(2000, 3, 26, 0, 0, 0, 0, time.UTC)
	tdate := time.Date(2030, 3, 26, 0, 0, 0, 0, time.UTC)
	params := Params{fdate: fdate, tdate: tdate, ratio: 100, minDuration: 30, maxDuration: 60}
	fsys := fstest.MapFS{
		"should_be_selected.mp4": {
			Data:    mocks.CreateData(50),
			Mode:    0755,
			ModTime: modTime,
		},
		"too_short.mp4": {
			Data:    mocks.CreateData(29),
			Mode:    0755,
			ModTime: modTime,
		},
		"too_long.mp4": {
			Data:    mocks.CreateData(65),
			Mode:    0755,
			ModTime: modTime,
		},
	}
	items, summary, _ := collectMediaContent("/home/Music", fsys, params)
	assert.Equal(t, "Should select one file", items[0].Name, "should_be_selected.mp4")
	assert.Equal(t, "Should select one file", summary.totalSelected, 1)
}

func TestCollectMediaContentDateFilter(t *testing.T) {
	fdate := time.Date(2022, 3, 26, 0, 0, 0, 0, time.UTC)
	tdate := time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC)
	params := Params{fdate: fdate, tdate: tdate, ratio: 100, minDuration: 0, maxDuration: math.MaxInt32}
	fsys := fstest.MapFS{
		"should_be_selected.mp4": {
			Data:    mocks.CreateData(120),
			Mode:    0755,
			ModTime: time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC),
		},
		"too_old.mp4": {
			Data:    mocks.CreateData(180),
			Mode:    0755,
			ModTime: time.Date(2022, 3, 25, 0, 0, 0, 0, time.UTC),
		},
		"too_recent.mp4": {
			Data:    mocks.CreateData(180),
			Mode:    0755,
			ModTime: time.Date(2024, 3, 27, 0, 0, 0, 0, time.UTC),
		},
	}
	items, summary, _ := collectMediaContent("/home/Music", fsys, params)
	assert.Equal(t, "Should select one file", items[0].Name, "should_be_selected.mp4")
	assert.Equal(t, "Should select one file", summary.totalSelected, 1)
}

func TestCollectMediaContentSkipFilter(t *testing.T) {
	modTime := time.Date(2020, 3, 26, 0, 0, 0, 0, time.UTC)
	fdate := time.Date(2000, 3, 26, 0, 0, 0, 0, time.UTC)
	tdate := time.Date(2030, 3, 26, 0, 0, 0, 0, time.UTC)
	params := Params{fdate: fdate, tdate: tdate, ratio: 100, minDuration: 0, maxDuration: math.MaxInt32, skipF: []string{"bad_artist"}}
	fsys := fstest.MapFS{
		"good_artist/should_be_selected.mp4": {
			Data:    mocks.CreateData(140),
			Mode:    0755,
			ModTime: modTime,
		},
		"bad_artist/ignored.mp4": {
			Data:    mocks.CreateData(150),
			Mode:    0755,
			ModTime: modTime,
		},
		"also_selected.mp4": {
			Data:    mocks.CreateData(160),
			Mode:    0755,
			ModTime: modTime,
		},
	}
	items, summary, _ := collectMediaContent("/home/Music", fsys, params)
	assert.Equal(t, "Id should be 0 for first item", items[0].Id, 0)
	assert.Equal(t, "Id should be 1 for second item", items[1].Id, 1)
	assert.Equal(t, "Should select two files", summary.totalSelected, 2)
}

func TestCollectMediaContentIncludeFilter(t *testing.T) {
	modTime := time.Date(2020, 3, 26, 0, 0, 0, 0, time.UTC)
	fdate := time.Date(2000, 3, 26, 0, 0, 0, 0, time.UTC)
	tdate := time.Date(2030, 3, 26, 0, 0, 0, 0, time.UTC)
	params := Params{fdate: fdate, tdate: tdate, ratio: 100, minDuration: 0, maxDuration: math.MaxInt32, includeF: []string{"good_artist"}}
	fsys := fstest.MapFS{
		"good_artist/should_be_selected.mp4": {
			Data:    mocks.CreateData(140),
			Mode:    0755,
			ModTime: modTime,
		},
		"bad_artist/ignored.mp4": {
			Data:    mocks.CreateData(150),
			Mode:    0755,
			ModTime: modTime,
		},
		"also_not_selected.mp4": {
			Data:    mocks.CreateData(160),
			Mode:    0755,
			ModTime: modTime,
		},
	}
	items, summary, _ := collectMediaContent("/home/Music", fsys, params)
	assert.Equal(t, "Should select one file", items[0].Name, "should_be_selected.mp4")
	assert.Equal(t, "Should select one file", summary.totalSelected, 1)
}

func TestCollectMediaContentSelector(t *testing.T) {
	fdate := time.Date(2000, 3, 26, 0, 0, 0, 0, time.UTC)
	tdate := time.Date(2030, 3, 26, 0, 0, 0, 0, time.UTC)
	params := Params{fdate: fdate, tdate: tdate, ratio: 0, minDuration: 0, maxDuration: math.MaxInt32}
	fsys := fstest.MapFS{
		"track_01.mp4": {
			Data:    mocks.CreateData(100),
			Mode:    0755,
			ModTime: time.Now(),
		},
		"track_02.mp4": {
			Data:    mocks.CreateData(100),
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	items, summary, _ := collectMediaContent("/home/Music", fsys, params)
	assert.Equal(t, "Should not select anything", len(items), 0)
	assert.Equal(t, "Should not select anything", summary.totalSelected, 0)
}

func _createMediaItems(length int) (items []MediaItem, indices []int) {
	for i := 0; i < length; {
		trackName := "track_" + strconv.Itoa(i) + ".mp4"
		items = append(items, MediaItem{Name: trackName, Id: i})
		indices = append(indices, i)
		i++
	}
	return
}

func _checkRandomness(original, randomized []int) (count int) {
	for i := range original {
		if original[i] == randomized[i] {
			count++
		}
	}
	return
}

func _getIndices(items []MediaItem) (newIndices []int) {
	for _, item := range items {
		newIndices = append(newIndices, item.Id)
	}
	return
}

func TestRandomizePlaylistWithoutStabilizer(t *testing.T) {
	mItems, indices := _createMediaItems(30)
	randomizePlaylist(mItems, len(indices)+1)

	newIndices := _getIndices(mItems)
	assert.NotEqualSlice(t, "Indices should be different", newIndices, indices)
}

func TestRandomizePlaylistWithStabilizer(t *testing.T) {
	mItems, indices := _createMediaItems(30)
	randomizePlaylist(mItems, 2)

	newIndices := _getIndices(mItems)
	assert.NotEqualSlice(t, "Indices should be different", newIndices, indices)
}

func TestRandomizePlaylistStabilizerLessRandom(t *testing.T) {
	nonStabilized, idxNonStabilized := _createMediaItems(100)
	randomizePlaylist(nonStabilized, len(idxNonStabilized)+1)
	newIndiecesNonStabilized := _getIndices(nonStabilized)

	stabilized, idxStabilized := _createMediaItems(100)
	randomizePlaylist(stabilized, 2)
	newIndicesStablized := _getIndices(stabilized)

	assert.NotEqualSlice(t, "Indices should be different", newIndiecesNonStabilized, idxNonStabilized)
	assert.NotEqualSlice(t, "Indices should be different", newIndicesStablized, idxStabilized)

	notMovedNonStabilized := _checkRandomness(idxNonStabilized, newIndiecesNonStabilized)
	notMovedStabilized := _checkRandomness(idxStabilized, newIndicesStablized)
	stabilizedLessRandom := notMovedStabilized > notMovedNonStabilized
	assert.Equal(t, "Randomness of not stabilized shuffle should be higher", stabilizedLessRandom, true)
}

func TestWritePlayList(t *testing.T) {
	var buf bytes.Buffer
	items := []MediaItem{
		{
			AbsPath:  "/home/Music/best track ever.mp4",
			Location: "/home/Music/best%20track%20ever.mp4",
			Name:     "track.mp4",
			Duration: 180,
			Id:       0,
		},
	}
	originalPl := buildPlayList(items, Options{Audio: true})
	writePlayList(originalPl, &buf)
	output := strings.Split(buf.String(), "\n")
	assert.Equal(t, "Output should be 14 rows", len(output), 14)
	assert.Equal(t, "Output should match", strings.TrimSpace(output[0]), "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[1]), "<playlist xmlns=\"http://xspf.org/ns/0/\" xmlns:vlc=\"http://www.videolan.org/vlc/playlist/ns/0/\" version=\"1\">")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[2]), "<title></title>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[3]), "<trackList>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[4]), "<track>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[5]), "<location>/home/Music/best%20track%20ever.mp4</location>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[6]), "<title>track.mp4</title>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[7]), "<duration>180</duration>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[8]), "<extension application=\"http://www.videolan.org/vlc/playlist/0\">")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[9]), "<vlc:id>0</vlc:id>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[10]), "</extension>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[11]), "</track>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[12]), "</trackList>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[13]), "</playlist>")
}

func TestWritePlayListWithNoAudioOption(t *testing.T) {
	var buf bytes.Buffer
	items := []MediaItem{
		{
			AbsPath:  "/home/Music/best track ever.mp4",
			Location: "/home/Music/best%20track%20ever.mp4",
			Name:     "track.mp4",
			Duration: 180,
			Id:       0,
		},
	}
	originalPl := buildPlayList(items, Options{Audio: false})
	writePlayList(originalPl, &buf)
	output := strings.Split(buf.String(), "\n")
	assert.Equal(t, "Output should be 15 rows", len(output), 15)
	assert.Equal(t, "Output should match", strings.TrimSpace(output[8]), "<extension application=\"http://www.videolan.org/vlc/playlist/0\">")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[9]), "<vlc:id>0</vlc:id>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[10]), "<vlc:option>no-audio</vlc:option>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[11]), "</extension>")
}

func TestWritePlayListWithStartOption(t *testing.T) {
	var buf bytes.Buffer
	items := []MediaItem{
		{
			AbsPath:  "/home/Music/best track ever.mp4",
			Location: "/home/Music/best%20track%20ever.mp4",
			Name:     "track.mp4",
			Duration: 180,
			Id:       0,
		},
	}
	originalPl := buildPlayList(items, Options{Audio: true, StartTime: 50})
	writePlayList(originalPl, &buf)
	output := strings.Split(buf.String(), "\n")
	assert.Equal(t, "Output should be 15 rows", len(output), 15)
	assert.Equal(t, "Output should match", strings.TrimSpace(output[8]), "<extension application=\"http://www.videolan.org/vlc/playlist/0\">")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[9]), "<vlc:id>0</vlc:id>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[10]), "<vlc:option>start-time=50</vlc:option>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[11]), "</extension>")
}

func TestWritePlayListWithEndOption(t *testing.T) {
	var buf bytes.Buffer
	items := []MediaItem{
		{
			AbsPath:  "/home/Music/best track ever.mp4",
			Location: "/home/Music/best%20track%20ever.mp4",
			Name:     "track.mp4",
			Duration: 180,
			Id:       0,
		},
	}
	originalPl := buildPlayList(items, Options{Audio: true, StopTime: 50})
	writePlayList(originalPl, &buf)
	output := strings.Split(buf.String(), "\n")
	assert.Equal(t, "Output should be 15 rows", len(output), 15)
	assert.Equal(t, "Output should match", strings.TrimSpace(output[8]), "<extension application=\"http://www.videolan.org/vlc/playlist/0\">")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[9]), "<vlc:id>0</vlc:id>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[10]), "<vlc:option>stop-time=50</vlc:option>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[11]), "</extension>")
}

func TestWritePlayListWithAllOption(t *testing.T) {
	var buf bytes.Buffer
	items := []MediaItem{
		{
			AbsPath:  "/home/Music/best track ever.mp4",
			Location: "/home/Music/best%20track%20ever.mp4",
			Name:     "track.mp4",
			Duration: 180,
			Id:       0,
		},
	}
	originalPl := buildPlayList(items, Options{Audio: false, StartTime: 30, StopTime: 50})
	writePlayList(originalPl, &buf)
	output := strings.Split(buf.String(), "\n")
	assert.Equal(t, "Output should be 15 rows", len(output), 17)
	assert.Equal(t, "Output should match", strings.TrimSpace(output[8]), "<extension application=\"http://www.videolan.org/vlc/playlist/0\">")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[9]), "<vlc:id>0</vlc:id>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[10]), "<vlc:option>no-audio</vlc:option>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[11]), "<vlc:option>start-time=30</vlc:option>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[12]), "<vlc:option>stop-time=50</vlc:option>")
	assert.Equal(t, "Output should match", strings.TrimSpace(output[13]), "</extension>")
}

func TestWritePlayListWriteError(t *testing.T) {
	err := writePlayList("", mocks.FakeWriter{})
	assert.ErrorRaised(t, "Error should be raised", err, true)
}

func TestBuildPlaylist(t *testing.T) {
	items := []MediaItem{
		{
			AbsPath:  "/home/Music/best track ever.mp4",
			Location: "/home/Music/best%20track%20ever.mp4",
			Name:     "track.mp4",
			Duration: 180,
			Id:       0,
		},
	}
	pl := buildPlayList(items, Options{Audio: false, StartTime: 100, StopTime: 180})
	assert.Equal(t, "Should have correct Xmlns value", pl.Xmlns, "http://xspf.org/ns/0/")
	assert.Equal(t, "Should have correct XmlnsVlc value", pl.XmlnsVlc, "http://www.videolan.org/vlc/playlist/ns/0/")
	assert.Equal(t, "Should have correct version", pl.Version, "1")

	assert.Equal(t, "Should have one track", len(pl.Tl.Tracks), 1)
	assert.Equal(t, "Should have correct extension application", pl.Tl.Tracks[0].Ext.Application, "http://www.videolan.org/vlc/playlist/0")
	assert.Equal(t, "Should have correct Id", pl.Tl.Tracks[0].Ext.Id, 0)
	assert.EqualSlice(t, "Should have all options added", pl.Tl.Tracks[0].Ext.Options, []string{"no-audio", "start-time=100", "stop-time=180"})
	assert.Equal(t, "Should have correct absolute path", pl.Tl.Tracks[0].Location, "/home/Music/best%20track%20ever.mp4")
	assert.Equal(t, "Should have correct title", pl.Tl.Tracks[0].Title, "track.mp4")
	assert.Equal(t, "Should have correct duration", pl.Tl.Tracks[0].Duration, 180)
}
