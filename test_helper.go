package gox

import "testing"

const (
	ErrorTextFormat = "%v not equal: [\nExpected: %v\nActual:   %v\n]\n"
)

func CompareString(label string, expected string, actual string, t *testing.T) bool {
	if expected != actual {
		t.Errorf(ErrorTextFormat, label, expected, actual)
		return true
	}

	return false
}

func CompareInt64(label string, expected int64, actual int64, t *testing.T) bool {
	if expected != actual {
		t.Errorf(ErrorTextFormat, label, expected, actual)
		return true
	}

	return false
}

func CompareInt(label string, expected int, actual int, t *testing.T) bool {
	if expected != actual {
		t.Errorf(ErrorTextFormat, label, expected, actual)
		return true
	}

	return false
}