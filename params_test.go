package main

import (
	"playmix/internal/assert"
	"playmix/internal/mocks"
	"testing"
	"testing/fstest"
	"time"
)

func TestParamsSetFileName(t *testing.T) {
	p := Params{}
	p.setFileName("myplaylist")
	assert.Equal(t, "Should set file name correctly", p.FileName, "myplaylist.xspf")
}

func TestParamsSetFileNameIfNotProvided(t *testing.T) {
	p := Params{}
	p.setFileName("")
	assert.Equal(t, "Should set file name correctly", p.FileName, "pl-test.xspf")
}

func TestParamsSetFileNameError(t *testing.T) {
	p := Params{}
	err := p.setFileName("bad_name.xspf")
	assert.ErrorRaised(t, "File param with extension specified raises error", err, true)
}

func TestParamsSetDateParams(t *testing.T) {
	p := Params{}
	p.setDateParams("20230326", "20240326")

	expectedFDate := time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC)
	expectedTDate := time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "Should convert string to fdate", expectedFDate, p.fdate)
	assert.Equal(t, "Should convert string to tdate", expectedTDate, p.tdate)
}

func TestParamsSetDateParamsError(t *testing.T) {
	p := Params{}
	err := p.setDateParams("invalid", "20240326")
	assert.ErrorRaised(t, "Setting date params should return error", err, true)
	err = p.setDateParams("20241212", "")
	assert.ErrorRaised(t, "Setting date params should return error", err, true)
	err = p.setDateParams("20241212", "20231212")
	assert.ErrorRaised(t, "Setting date params should return error", err, true)
}

func TestParseOptFile(t *testing.T) {
	p := Params{}
	data := []byte(`{ "media_path": "path/to/media/", "marquee": {"text": "test", "color":"black", "opacity":50, "position":"center"}}`)
	fn := "options.json"
	f := fstest.MapFS{
		fn: {
			Data:    data,
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	err := p.parseOptFile(f, fn)
	assert.Equal(t, "Parsing should succeed", err, nil)
	expectedMarquee := Marquee{
		Text:     "test",
		Color:    "black",
		Opacity:  50,
		Position: "center",
	}
	assert.Equal(t, "Marquee should be set", p.MarqueeOptions, expectedMarquee)
}

func TestParseOptFileOpenError(t *testing.T) {
	p := Params{}
	err := p.parseOptFile(fstest.MapFS{}, "no_such_file.json")
	assert.ErrorRaised(t, "Invalid file name should raise error", err, true)
}

func TestParseOptFileNoFileRaisesError(t *testing.T) {
	p := Params{}
	err := p.parseOptFile(fstest.MapFS{}, "")
	assert.ErrorRaised(t, "No file name provided should raise error", err, true)
}

func TestParseOptFileReadAllError(t *testing.T) {
	p := Params{}
	fsys := mocks.FakeSys{}
	err := p.parseOptFile(fsys, "readMock.json")
	assert.ErrorRaised(t, "Invalid file name should raise error", err, true)
}

func TestParseOptUnMarshallError(t *testing.T) {
	p := Params{}
	data := []byte("\\")
	fn := "options.json"
	f := fstest.MapFS{
		fn: {
			Data:    data,
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	err := p.parseOptFile(f, fn)
	assert.ErrorRaised(t, "Should raise unmarshall error", err, true)
}

func TestParseOptFileInvalidPath(t *testing.T) {
	p := Params{}
	data := []byte(`{"media_path":""}`)
	fn := "options.json"
	f := fstest.MapFS{
		fn: {
			Data:    data,
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	err := p.parseOptFile(f, fn)
	assert.ErrorRaised(t, "Should raise media_path needs to be set error", err, true)
}

func TestParseOptFileInvalidPosition(t *testing.T) {
	p := Params{}
	data := []byte(`{"media_path":"/media/", "marquee": {"text": "test", "color": "blue", "opacity":100, "position":"invalid"}}`)
	fn := "options.json"
	f := fstest.MapFS{
		fn: {
			Data:    data,
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	err := p.parseOptFile(f, fn)
	assert.ErrorRaised(t, "Should raise unrecognized position error", err, true)
}

func TestParseOptFileInvalidColor(t *testing.T) {
	p := Params{}
	data := []byte(`{"media_path":"/media/", "marquee": {"text": "test", "color": "invalid", "opacity":100, "position":"center"}}`)
	fn := "options.json"
	f := fstest.MapFS{
		fn: {
			Data:    data,
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	err := p.parseOptFile(f, fn)
	assert.ErrorRaised(t, "Should raise unrecognized color error", err, true)
}

func TestParseOptFileInvalidTimes(t *testing.T) {
	p := Params{}
	data := []byte(`{"media_path":"/media/", "play_options": {"start_time": 100, "stop_time": 10}}`)
	fn := "options.json"
	f := fstest.MapFS{
		fn: {
			Data:    data,
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	err := p.parseOptFile(f, fn)
	assert.ErrorRaised(t, "Should raise invalid time error", err, true)
}

func TestParseOptFileInvalidRatio(t *testing.T) {
	p := Params{}
	data := []byte(`{"media_path":"/media/", "randomizer_options": {"ratio": 110 }}`)
	fn := "options.json"
	f := fstest.MapFS{
		fn: {
			Data:    data,
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	err := p.parseOptFile(f, fn)
	assert.ErrorRaised(t, "Should raise invalid time error", err, true)
}

func TestParseOptInvalidFilterOptions(t *testing.T) {
	p := Params{}
	data := []byte(`{"media_path":"/media/", "filter_options": {"include_folder": ["folderA"], "skip_folder": ["folderB"]}}`)
	fn := "options.json"
	f := fstest.MapFS{
		fn: {
			Data:    data,
			Mode:    0755,
			ModTime: time.Now(),
		},
	}
	err := p.parseOptFile(f, fn)
	assert.ErrorRaised(t, "Should raise invalid filter options error", err, true)
}
