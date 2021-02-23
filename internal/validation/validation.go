package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Value return nil if val is equal (reflect.DeepEqual) to any of the values in validVals
func Value(val interface{}, validVals ...interface{}) error {
	for _, t := range validVals {
		if reflect.DeepEqual(val, t) {
			return nil
		}
	}
	var sValidVals []string
	for _, t := range validVals {
		sValidVals = append(sValidVals, fmt.Sprintf("%v", t))
	}
	return fmt.Errorf("%q is not any of %s", val, strings.Join(sValidVals, ", "))
}

// UUID4 return nil if val is a valid uuid
func UUID4(val string) error {
	b := []byte(strings.ToLower(val))
	if len(b) != 36 {
		return fmt.Errorf("uuid4: length is not 36")
	}
	for pos, c := range b {
		switch pos {
		case 8, 13, 18, 23:
			if c != '-' {
				return fmt.Errorf("uuid4: pos %d invalid delimiter character %c", pos, c)
			}
		case 14:
			if c != '4' {
				return fmt.Errorf("uuid4: pos %d invalid version %c", pos, c)
			}
		default:
			if (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') || (c >= '0' && c <= '9') {
				continue
			}
			return fmt.Errorf("uuid4: pos %d invalid character %c", pos, c)
		}
	}
	return nil
}

// Numeric returns nil if v is of a numeric type or a string that can be parsed as a number.
func Numeric(v interface{}) error {
	switch v.(type) {
	case int, uint, int32, uint32, int64, uint64, float32, float64:
		return nil
	default:
		if _, err := strconv.ParseFloat(fmt.Sprintf("%s", v), 64); err == nil {
			return nil
		}
	}
	return fmt.Errorf("value %q is not numeric", v)
}
