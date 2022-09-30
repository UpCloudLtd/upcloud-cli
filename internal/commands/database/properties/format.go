package databaseproperties

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

func formatAlternatives(val interface{}) (text.Colors, string, error) {
	if val == nil {
		return nil, "", nil
	}

	if stringVal, ok := val.(string); ok {
		return nil, stringVal, nil
	}

	if ifaceSliceVal, ok := val.([]interface{}); ok {
		return nil, alternativesString(ifaceSliceVal), nil
	}

	return nil, fmt.Sprintf("%+v", val), nil
}

func maxStringLen(strings []string) int {
	max := 0
	for _, str := range strings {
		if strLen := len(str); strLen > max {
			max = strLen
		}
	}
	return max
}

func alternativesString(values []interface{}) string {
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		return fmt.Sprintf("%+v", values[0])
	}

	strs := make([]string, len(values))
	for i, value := range values {
		strs[i] = fmt.Sprintf("%+v", value)
	}

	whitespace := " "
	if maxStringLen(strs) > 15 {
		whitespace = "\n"
	}

	str := strings.Join(strs[:len(strs)-1], ","+whitespace)
	return str + fmt.Sprintf(" or%s%s", whitespace, strs[len(values)-1])
}
