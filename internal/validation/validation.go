package validation

import (
	"fmt"
	"reflect"
	"strings"
)

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

func Uuid4(val string) error {
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
