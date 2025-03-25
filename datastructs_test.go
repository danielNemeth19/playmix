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
