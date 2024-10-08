package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcatStringsSingleString(t *testing.T) {
	result := ConcatStrings("hello")
	expected := "hello"
	assert.Equal(t, expected, result)
}

func TestConcatStringsMultipleStrings(t *testing.T) {
	result := ConcatStrings("hello", "world", "foo", "bar")
	expected := "hello/world/foo/bar"
	assert.Equal(t, expected, result)
}

func TestConcatStringsWithEmptyStrings(t *testing.T) {
	result := ConcatStrings("hello", "", "world", "", "foo", "bar")
	expected := "hello/world/foo/bar"
	assert.Equal(t, expected, result)
}

func TestConcatStringsAllEmptyStrings(t *testing.T) {
	result := ConcatStrings("", "", "")
	expected := ""
	assert.Equal(t, expected, result)
}

func TestConcatStringsEmptyInput(t *testing.T) {
	result := ConcatStrings()
	expected := ""
	assert.Equal(t, expected, result)
}
