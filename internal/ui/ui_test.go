package ui

import (
	"testing"
)

func TestConcatStringsSingleString(t *testing.T) {
	result := ConcatStrings("hello")
	expected := "hello"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestConcatStringsMultipleStrings(t *testing.T) {
	result := ConcatStrings("hello", "world", "foo", "bar")
	expected := "hello/world/foo/bar"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestConcatStringsWithEmptyStrings(t *testing.T) {
	result := ConcatStrings("hello", "", "world", "", "foo", "bar")
	expected := "hello/world/foo/bar"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestConcatStringsAllEmptyStrings(t *testing.T) {
	result := ConcatStrings("", "", "")
	expected := ""
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestConcatStringsEmptyInput(t *testing.T) {
	result := ConcatStrings()
	expected := ""
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
