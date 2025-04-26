package main

import (
	"fmt"
	"playmix/internal/assert"
	"testing"
)

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
		StopTime: 50,
	}
	err := opts.ValidateTimes()
	assert.ErrorRaised(t, "Start time should not be smaller than stop time", err, true)
}

func TestValidateTimesNotProvided(t *testing.T) {
	opts := PlayOptions{
		StartTime: 0,
		StopTime: 0,
	}
	err := opts.ValidateTimes()
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

