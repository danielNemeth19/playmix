package main

import (
	"fmt"
	"playmix/internal/assert"
	"testing"
)

func TestValidatePath(t *testing.T) {
	fOpts := FileOptions{MediaPath:  "/home/user/Music/"}
	err := fOpts.validatePath()
	expected := "/home/user/Music/"
	assert.Equal(t, "Should parse root path", fOpts.MediaPath, expected)
	assert.ErrorRaised(t, "No error", err, false)
}

func TestValidatePathNormalized(t *testing.T) {
	fOpts := FileOptions{MediaPath:  "/home/user/Music"}
	err := fOpts.validatePath()
	expected := "/home/user/Music/"
	assert.Equal(t, "Should parse root path normalized", fOpts.MediaPath, expected)
	assert.ErrorRaised(t, "No error", err, false)
}

func TestValidatePathRaisesError(t *testing.T) {
	fOpts := FileOptions{MediaPath:  ""}
	err := fOpts.validatePath()
	assert.ErrorRaised(t, "Should raise error", err, true)
}

func TestValidateColor(t *testing.T) {
	m := Marquee{Color: "red"}
	err := m.validateColor()
	assert.ErrorRaised(t, "Error should not be raised for valid color", err, false)
}

func TestValidateColorError(t *testing.T) {
	m := Marquee{Color: "invalid"}
	err := m.validateColor()
	assert.ErrorRaised(t, "Error should be raised for invalid color", err, true)
}

func TestValidatePosition(t *testing.T) {
	m := Marquee{Position: "center"}
	err := m.validatePosition()
	assert.ErrorRaised(t, "Error should not be raised for valid position", err, false)
}

func TestValidatePositionError(t *testing.T) {
	m := Marquee{Position: "invalid"}
	err := m.validatePosition()
	assert.ErrorRaised(t, "Error should be raised for invalid position", err, true)
}

func TestRemapColor(t *testing.T) {
	for color, colorRGBA := range colorMap {
		m := Marquee{
			Text:     "Test",
			Color:    color,
			Opacity:  100,
			Position: "center",
		}
		res := m.remapColor()
		message := fmt.Sprintf("Should map %s to %d\n", color, colorRGBA)
		assert.Equal(t, message, res, colorRGBA)
	}
}

func TestRemapColorDefault(t *testing.T) {
	m := Marquee{
		Text:     "Test",
		Color:    "invalid",
		Opacity:  100,
		Position: "center",
	}
	colorRGBA := m.remapColor()
	assert.Equal(t, "Should default to red", colorRGBA, colorMap["red"])
}

func TestRemapPosition(t *testing.T) {
	for pos, vlcPos := range textPositionMap {
		m := Marquee{
			Text:     "Test",
			Opacity:  100,
			Position: pos,
		}
		res := m.remapPosition()
		message := fmt.Sprintf("Should map %s to %d\n", pos, vlcPos)
		assert.Equal(t, message, res, vlcPos)
	}
}

func TestRemapPositionDefault(t *testing.T) {
	m := Marquee{
		Text:     "Test",
		Color:    "cyan",
		Opacity:  100,
		Position: "invalid",
	}
	vlcPos := m.remapPosition()
	assert.Equal(t, "Should default to disable", vlcPos, textPositionMap["disable"])
}

func TestValidateTimes(t *testing.T) {
	opts := PlayOptions{
		StartTime: 100,
		StopTime:  50,
	}
	err := opts.validateTimes()
	assert.ErrorRaised(t, "Start time should not be smaller than stop time", err, true)
}

func TestValidateTimesNotProvided(t *testing.T) {
	opts := PlayOptions{
		StartTime: 0,
		StopTime:  0,
	}
	err := opts.validateTimes()
	assert.ErrorRaised(t, "Should not raise error", err, false)
}

func TestStringifyAudio(t *testing.T) {
	opts := PlayOptions{
		Audio: true,
	}
	xmlValue := opts.StringifyAudio()
	assert.Equal(t, "Should be emtpy string", xmlValue, "")
}

func TestStringifyNoAudio(t *testing.T) {
	opts := PlayOptions{
		Audio: false,
	}
	xmlValue := opts.StringifyAudio()
	assert.Equal(t, "Should be emtpy string", xmlValue, "no-audio")
}

func TestStringifyStartTime(t *testing.T) {
	opts := PlayOptions{
		StartTime: 100,
	}
	xmlValue := opts.StringifyStartTime()
	assert.Equal(t, "StartTime should be start-time=100 in xml", xmlValue, "start-time=100")
}

func TestStringifyStopTime(t *testing.T) {
	opts := PlayOptions{
		StopTime: 200,
	}
	xmlValue := opts.StringifyStopTime()
	assert.Equal(t, "StopTime should be stop-time=200 in xml", xmlValue, "stop-time=200")
}

func TestParamsValidateRatio(t *testing.T) {
	opts := RandomizerOptions{
		Ratio: 100,
	}
	err := opts.validateRatio()
	assert.Equal(t, "Validate return nil", nil, err)

	opts = RandomizerOptions{
		Ratio: 150,
	}
	err = opts.validateRatio()
	assert.ErrorRaised(t, "Validate ratio should return error", err, true)
}
