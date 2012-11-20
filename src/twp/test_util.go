package twp

import (
    "testing"
)

func pError(err error, t *testing.T) {
    t.Errorf("Error: %s\n", err)
}

func pFatal(err error, t *testing.T) {
    t.Fatalf("Fatal: %s\n", err)
}

func verify(a, b interface{}, t *testing.T) {
    if a != b {
        t.Fatalf("Expected %s to be %s.\n", a, b)
    }
}