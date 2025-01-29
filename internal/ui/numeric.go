package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type spec struct {
	th     uint
	suffix string
}

var siMultiples = [...]spec{
	{th: 1 * 1000 * 1000 * 1000 * 1000, suffix: "T"},
	{th: 1 * 1000 * 1000 * 1000, suffix: "G"},
	{th: 1 * 1000 * 1000, suffix: "M"},
	{th: 1 * 1000, suffix: "K"},
	{th: 1, suffix: ""},
}

var binaryMultiples = [...]spec{
	{th: 1 * 1024 * 1024 * 1024 * 1024, suffix: "Ti"},
	{th: 1 * 1024 * 1024 * 1024, suffix: "Gi"},
	{th: 1 * 1024 * 1024, suffix: "Mi"},
	{th: 1 * 1024, suffix: "Ki"},
	{th: 1, suffix: ""},
}

func abbrevInt(raw uint, specs [5]spec) string {
	var (
		val    float64
		suffix string
	)
	for _, spec := range specs {
		if spec.th > raw {
			continue
		}
		if float64(raw)/float64(spec.th) < 0 {
			continue
		}
		val = float64(raw) / float64(spec.th)
		suffix = spec.suffix
		break
	}
	if _, frac := math.Modf(val); frac > 0 {
		return fmt.Sprintf("%.2f%s", val, suffix)
	}
	return fmt.Sprintf("%.0f%s", val, suffix)
}

func parseAbbrevInt(s string, specs [5]spec) (uint, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("invalid value %q", s)
	}
	for _, spec := range specs {
		if strings.HasSuffix(s, spec.suffix) {
			if spec.suffix == "" && !unicode.IsDigit(rune(s[len(s)-1])) {
				s = s[0 : len(s)-1]
			}
			v, err := strconv.ParseFloat(strings.TrimRight(s, spec.suffix), 64)
			if err != nil {
				return 0, err
			}
			return uint(v * float64(spec.th)), nil
		}
	}
	return 0, fmt.Errorf("invalid value %q", s)
}

// AbbrevNum returns a string with the given number abbreviated with SI formatting (eg 1000 = 1k)
func AbbrevNum(raw uint) string {
	return abbrevInt(raw, siMultiples)
}

// AbbrevNumBinaryPrefix returns a string with the given number abbreviated with binary formatting (eg 1024 = 1ki)
func AbbrevNumBinaryPrefix(raw uint) string {
	return abbrevInt(raw, binaryMultiples)
}

// FormatBytes returns a string with the given number interpreted as bytes and abbreviated with binary formatting (eg 1024 = 1KiB)
func FormatBytes(n int) string {
	if n < 0 || n > math.MaxUint32 {
		return fmt.Sprintf("%sB", "-1")
	}

	return fmt.Sprintf("%sB", AbbrevNumBinaryPrefix(uint(n)))
}

// ParseAbbrevNum parses a string formatted to a uint in SI unit style. (eg. "1k" = 1000)
func ParseAbbrevNum(s string) (uint, error) {
	return parseAbbrevInt(s, siMultiples)
}

// ParseAbbrevNumBinaryPrefix parses a string formatted to a uint in binary units. (eg. "1Ki" = 1000)
func ParseAbbrevNumBinaryPrefix(s string) (uint, error) {
	return parseAbbrevInt(s, binaryMultiples)
}
