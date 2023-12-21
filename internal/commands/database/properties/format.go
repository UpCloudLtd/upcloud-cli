package databaseproperties

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

func formatAlternatives(val interface{}) (text.Colors, string, error) {
	return formatStringSlice(val, "or")
}

func formatProperties(val interface{}) (text.Colors, string, error) {
	return formatStringSlice(val, "and")
}

func formatStringSlice(val interface{}, andOrOr string) (text.Colors, string, error) {
	if val == nil {
		return nil, "", nil
	}

	if stringVal, ok := val.(string); ok {
		return nil, stringVal, nil
	}

	if ifaceSliceVal, ok := val.([]interface{}); ok {
		return nil, alternativesString(ifaceSliceVal, andOrOr), nil
	}

	if stringSliceVal, ok := val.([]string); ok {
		ifaceSliceVal := []interface{}{}
		for _, i := range stringSliceVal {
			ifaceSliceVal = append(ifaceSliceVal, i)
		}
		return nil, alternativesString(ifaceSliceVal, andOrOr), nil
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

func alternativesString(values []interface{}, andOrOr string) string {
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
	if maxStringLen(strs) > 15 || len(strs) > 3 {
		whitespace = "\n"
	}

	str := strings.Join(strs[:len(strs)-1], ","+whitespace)
	return str + fmt.Sprintf(" %s%s%s", andOrOr, whitespace, strs[len(values)-1])
}
