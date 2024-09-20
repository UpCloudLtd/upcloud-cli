package format

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

func StringSliceOr(val interface{}) (text.Colors, string, error) {
	return formatStringSlice(val, "or", false)
}

func StringSliceAnd(val interface{}) (text.Colors, string, error) {
	return formatStringSlice(val, "and", false)
}

func StringSliceSingleLineAnd(val interface{}) (text.Colors, string, error) {
	return formatStringSlice(val, "and", true)
}

func formatStringSlice(val interface{}, andOrOr string, singleLine bool) (text.Colors, string, error) {
	if val == nil {
		return nil, "", nil
	}

	if stringVal, ok := val.(string); ok {
		return nil, stringVal, nil
	}

	if ifaceSliceVal, ok := toIfaceSlice(val); ok {
		return nil, stringSliceString(ifaceSliceVal, andOrOr, singleLine), nil
	}

	return nil, fmt.Sprintf("%+v", val), nil
}

func toIfaceSlice(val interface{}) ([]interface{}, bool) {
	if reflect.TypeOf(val).Kind() == reflect.Slice {
		ifaceSliceVal := []interface{}{}
		reflectedVal := reflect.ValueOf(val)
		for i := 0; i < reflectedVal.Len(); i++ {
			ifaceSliceVal = append(ifaceSliceVal, reflectedVal.Index(i).Interface())
		}
		return ifaceSliceVal, true
	}
	return nil, false
}

func maxStringLen(strings []string) int {
	maxLen := 0
	for _, str := range strings {
		if strLen := len(str); strLen > maxLen {
			maxLen = strLen
		}
	}
	return maxLen
}

func stringSliceString(values []interface{}, andOrOr string, singleLine bool) string {
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
	if !singleLine && (maxStringLen(strs) > 15 || len(strs) > 3) {
		whitespace = "\n"
	}

	str := strings.Join(strs[:len(strs)-1], text.FgHiBlack.Sprint(",")+whitespace)
	return str + fmt.Sprintf(" %s%s%s", text.FgHiBlack.Sprint(andOrOr), whitespace, strs[len(values)-1])
}
