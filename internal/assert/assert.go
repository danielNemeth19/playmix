package assert

import (
    "testing"
    "slices"
)

func Equal[T comparable](t *testing.T, name string, actual, expected T){
    t.Helper()
    if actual != expected {
		t.Errorf("Failed %s: got %v, expected %v\n", name, actual, expected)
    }
}

func ErrorRaised(t *testing.T, name string, err error, expected bool) {
    t.Helper()
    if expected && err == nil {
        t.Errorf("Failed %s: expected error, got: %v\n", name, err)
    } else if !expected && err != nil {
        t.Errorf("Failed %s: not expected an error, got: %v\n", name, err)
    }
}

func EqualSlice[T comparable](t *testing.T, name string, actual, expected []T){
    t.Helper()
    if !slices.Equal(actual, expected) {
		t.Errorf("Failed %s: got %v, expected %v\n", name, actual, expected)
    }
}

func NotEqualSlice[T comparable](t *testing.T, name string, actual, expected []T){
    t.Helper()
    if slices.Equal(actual, expected) {
		t.Errorf("Failed %s: got %v, expected %v\n", name, actual, expected)
    }
}
