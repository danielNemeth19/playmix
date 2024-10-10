package main

import (
	"playmix/internal/assert"
	"testing"
	"time"
)

func TestParamsValidateRatio(t *testing.T) {
	err := validateRatio(100)
	assert.Equal(t, "Validate return nil", nil, err)

	err = validateRatio(-50)
	assert.ErrorRaised(t, "Validate ratio should return error", err, true)

	err = validateRatio(150)
	assert.ErrorRaised(t, "Validate ratio should return error", err, true)
}

func TestParamsParseNil(t *testing.T) {
	input := ""
	expected := []string{}
	parsed := parseFolder(input)
	assert.EqualSlice(t, "Should return empty slice", parsed, expected)
}

func TestParamsParseOneItem(t *testing.T) {
	input := "abc"
	expected := []string{"abc"}
	parsed := parseFolder(input)
	assert.EqualSlice(t, "Should return slice", parsed, expected)
}

func TestParamsParse(t *testing.T) {
	input := "a,b,c"
	expected := []string{"a", "b", "c"}
	parsed := parseFolder(input)
	assert.EqualSlice(t, "Should return multiple folders in slice", parsed, expected)
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
}

func TestParamsSetFolderParams(t *testing.T) {
	p := Params{}
	p.setFolderParams("folderA,folderB", "")

	expectedIncludeF := []string{"folderA", "folderB"}
	expectedSkipF := []string{}
	assert.EqualSlice(t, "Should parse include folder param to slice", p.includeF, expectedIncludeF)
	assert.EqualSlice(t, "Should parse skip folder param to slice", p.skipF, expectedSkipF)
}

func TestParamsSetFolderParamsError(t *testing.T) {
	p := Params{}
	err := p.setFolderParams("folderA,folderB", "folderC,folderD")
	assert.ErrorRaised(t, "Folder params are mutually exclusive", err, true)
}
