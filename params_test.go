package main

import (
	"playmix/internal/assert"
	"slices"
	"testing"
)

func TestParamsParse(t *testing.T) {
	input := "a,b,c"
	expected := []string{"a", "b", "c"}
	parsed := parse(input)
	assert.Equal(t, "Parse should return slice", len(expected), 3)
	if !slices.Equal(parsed, expected) {
		t.Errorf("Parse should return slice: %v, got %v\n", expected, parsed)
	}
}

func TestParamsParseOneItem(t *testing.T) {
	input := "d"
	expected := []string{"d"}
	parsed := parse(input)
	assert.Equal(t, "Parse should return slice", len(expected), 1)
	if !slices.Equal(parsed, expected) {
		t.Errorf("Parse should return slice: %v, got %v\n", expected, parsed)
	}
}

func TestParamsParseNil(t *testing.T) {
	input := ""
	expected := []string{}
	parsed := parse(input)
	assert.Equal(t, "Parse should return empty slice", len(expected), 0)
	if !slices.Equal(parsed, expected) {
		t.Errorf("Parse should return empy slice: %v, got %v\n", expected, parsed)
	}
}
