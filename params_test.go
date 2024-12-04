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
	parsed := parseParam(input)
	assert.EqualSlice(t, "Should return empty slice", parsed, expected)
}

func TestParamsParseOneItem(t *testing.T) {
	input := "abc"
	expected := []string{"abc"}
	parsed := parseParam(input)
	assert.EqualSlice(t, "Should return slice", parsed, expected)
}

func TestParamsParse(t *testing.T) {
	input := "a,b,c"
	expected := []string{"a", "b", "c"}
	parsed := parseParam(input)
	assert.EqualSlice(t, "Should return multiple folders in slice", parsed, expected)
}

func TestParamsSetFileName(t *testing.T) {
	p := Params{}
	p.setFileName("myplaylist")
	assert.Equal(t, "Should set file name correctly", p.fileName, "myplaylist.xspf")
}

func TestParamsSetFileNameIfNotProvided(t *testing.T) {
	p := Params{}
	p.setFileName("")
	assert.Equal(t, "Should set file name correctly", p.fileName, "pl-test.xspf")
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

func TestSetOptionsNoAudio(t *testing.T) {
	p := Params{}
	p.setOptions("no-audio")
	assert.Equal(t, "Options for audio should be set", p.options.audio, "no-audio")
}

func TestSetOptionsStartTime(t *testing.T) {
	p := Params{}
	err := p.setOptions("start-time=60")
	assert.Equal(t, "no error raised", err, nil)
	assert.Equal(t, "Option for start-time should be set", p.options.StartTime, 60)
}

func TestSetOptionsStartTimeRaisesError(t *testing.T) {
	p := Params{}
	err := p.setOptions("start-time=-60")
	assert.ErrorRaised(t, "Negative value should raise an error", err, true)
}

func TestSetOptionsEndTime(t *testing.T) {
	p := Params{}
	err := p.setOptions("end-time=120")
	assert.Equal(t, "no error raised", err, nil)
	assert.Equal(t, "Option for end-time should be set", p.options.EndTime, 120)
}

func TestSetOptionsEndTimeRaisesError(t *testing.T) {
	p := Params{}
	err := p.setOptions("end-time=-120")
	assert.ErrorRaised(t, "Negative value should raise an error", err, true)
}

// TODO:Missing reflection error tests

func TestSetOptionsNoOptions(t *testing.T) {
	p := Params{}
	p.setOptions("")
	assert.Equal(t, "If no option struct should be empty", p.options, Options{})
}
