package main

import (
	"fmt"
	"playmix/internal/assert"
	"slices"
	"testing"
    "time"
)

func TestParamsDateStrings(t *testing.T) {
	params := Params{
		fdate: "20230326",
		tdate: "20240326",
	}
    expectedFDate := time.Date(2023, 3, 26, 0, 0, 0, 0, time.UTC)
    got := params.parseDateString("fdate")
    assert.Equal(t, "Should convert string to date", expectedFDate, got)
    
    expectedTDate := time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC)
    got = params.parseDateString("tdate")
    assert.Equal(t, "Should convert string to date", expectedTDate, got)
}

func TestParamsDateWrongKey(t *testing.T) {
    params := Params {
        fdate : "20220101",
        tdate : "20220201",
    }
    got := params.parseDateString("invalidkey")
    fmt.Println(got)
}

func TestParamsValidate(t *testing.T) {
	err := validateParams("", "", 100)
	assert.Equal(t, "Validate return nil\n", nil, err)
	err = validateParams("folder1", "", 100)
	assert.Equal(t, "Validate return nil\n", nil, err)
	err = validateParams("folderA", "folderB", 100)
	if err == nil {
		t.Errorf("Validate should return error, got: %v\n", err)
	}
	err = validateParams("", "", 150)
	fmt.Println(err)
	if err == nil {
		t.Errorf("Validate should return error, got: %v\n", err)
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

func TestParamsParseOneItem(t *testing.T) {
	input := "abc"
	expected := []string{"abc"}
	parsed := parse(input)
	assert.Equal(t, "Parse should return slice", len(expected), 1)
	if !slices.Equal(parsed, expected) {
		t.Errorf("Parse should return slice: %v, got %v\n", expected, parsed)
	}
}

func TestParamsParse(t *testing.T) {
	input := "a,b,c"
	expected := []string{"a", "b", "c"}
	parsed := parse(input)
	assert.Equal(t, "Parse should return slice", len(expected), 3)
	if !slices.Equal(parsed, expected) {
		t.Errorf("Parse should return slice: %v, got %v\n", expected, parsed)
	}
}
