package assert

import (
    "testing"
)

func Equal[T comparable](t *testing.T, name string, actual, expected T){
    t.Helper()
    if actual != expected {
		t.Errorf("Failed %s: got %v, expected %v\n", name, actual, expected)
    }
}
