package main

import (
	"playmix/internal/assert"
	"slices"
	"testing"
    "time"
)

func TestParamsValidateRatio(t *testing.T) {
	err := validateRatio(100)
	assert.Equal(t, "Validate return nil\n", nil, err)
	err = validateRatio(-50)
	if err == nil {
		t.Errorf("Validate ratio should return error, got: %v\n", err)
	}
	err = validateRatio(150)
	if err == nil {
		t.Errorf("Validate ratio should return error, got: %v\n", err)
	}
}

func TestParamsParseNil(t *testing.T) {
	input := ""
	expected := []string{}
	parsed := parseFolder(input)
	assert.Equal(t, "Parse should return empty slice", len(expected), 0)
	if !slices.Equal(parsed, expected) {
		t.Errorf("Parse should return empy slice: %v, got %v\n", expected, parsed)
	}
}

func TestParamsParseOneItem(t *testing.T) {
	input := "abc"
	expected := []string{"abc"}
	parsed := parseFolder(input)
	assert.Equal(t, "Parse should return slice", len(expected), 1)
	if !slices.Equal(parsed, expected) {
		t.Errorf("Parse should return slice: %v, got %v\n", expected, parsed)
	}
}

func TestParamsParse(t *testing.T) {
	input := "a,b,c"
	expected := []string{"a", "b", "c"}
	parsed := parseFolder(input)
	assert.Equal(t, "Parse should return slice", len(expected), 3)
	if !slices.Equal(parsed, expected) {
		t.Errorf("Parse should return slice: %v, got %v\n", expected, parsed)
	}
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
    if err == nil {
        t.Errorf("Setting date params should return error, got: %v\n", err)
    }
    err = p.setDateParams("20241212", "")
    if err == nil {
        t.Errorf("Setting date params should return error, got: %v\n", err)
    }
}

func TestParamsSetFolderParams(t *testing.T) {
	p := Params{}
    p.setFolderParams("folderA,folderB", "")

    expectedIncludeF := []string{"folderA", "folderB"}
    expectedSkipF := []string{}
	if !slices.Equal(expectedIncludeF, p.includeF) {
		t.Errorf("SetFolderParams should return slice: %v, got %v\n", expectedIncludeF, p.includeF)
	}
	if !slices.Equal(expectedSkipF, p.skipF) {
		t.Errorf("SetFolderParams should return slice: %v, got %v\n", expectedSkipF, p.skipF)
	}
}

func TestParamsSetFolderParamsError(t *testing.T) {
	p := Params{}
    err := p.setDateParams("folderA,folderB", "folderC,folderD")
    if err == nil {
        t.Errorf("Setting folder params should return error, got: %v\n", err)
    }
}

