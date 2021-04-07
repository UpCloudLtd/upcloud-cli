package commands

import (
	"encoding/csv"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/validation"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"reflect"
	"strings"
)

// Parse parses a complex, querystring-type argument from in and returns all the parts found
// eg. `--foo bar=baz,flop=flip` returns `[]string{"bar","baz","flop","flip"}`
func Parse(in string) ([]string, error) {
	var result []string
	reader := csv.NewReader(strings.NewReader(in))
	args, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for _, arg := range args {
		result = append(result, strings.Split("--"+arg, "=")...)
	}
	return result, nil
}

// ToArray turns an interface{} to a slice of interface{}s.
// If the underlying type is also a slice, the elements will be returned as the return values elements..
// Otherwise, the input element is wrapped in a slice.
func ToArray(in interface{}) []interface{} {
	var elems []interface{}
	if reflect.TypeOf(in).Kind() == reflect.Slice {
		is := reflect.ValueOf(in)
		for i := 0; i < is.Len(); i++ {
			elems = append(elems, is.Index(i).Interface())
		}
	} else {
		elems = append(elems, in)
	}
	return elems
}

// SearchResources is a convenience method to map a list of resources to uuids.
// Any input strings that are uuids are returned as such and any other string is
// passed on to searchFn, the results of which are passed on to getUUID which is
// expected to return a uuid.
func SearchResources(
	ids []string,
	searchFn func(id string) (interface{}, error),
	getUUID func(interface{}) string,
) ([]string, error) {
	var result []string
	for _, id := range ids {
		if err := validation.UUID4(id); err == nil {
			result = append(result, id)
		} else {
			matchedResults, err := searchFn(id)
			if err != nil {
				return nil, err
			}

			for _, resource := range ToArray(matchedResults) {
				result = append(result, getUUID(resource))
			}
		}
	}
	return result, nil
}

// StateColour is a helper mapping states to colors
func StateColour(state string) text.Colors {
	switch state {
	case upcloud.ServerStateStarted:
		return text.Colors{text.FgGreen}
	case upcloud.ServerStateError:
		return text.Colors{text.FgHiRed, text.Bold}
	case upcloud.ServerStateMaintenance:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// BoolFromString parses a string and returns *upcloud.Boolean
func BoolFromString(b string) (*upcloud.Boolean, error) {
	// TODO: why does this return a pointer? this should (eventually) not be needed as tristate flags
	// should be handled much more easily than with this approach
	var result upcloud.Boolean
	switch b {
	case "true":
		result = upcloud.FromBool(true)
	case "false":
		result = upcloud.FromBool(false)
	default:
		return nil, fmt.Errorf("invalid boolean value %s", b)
	}
	return &result, nil
}
