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

func SearchResources(
	ids []string,
	searchFn func(id string) (interface{}, error),
	getUUID func(interface{}) string,
) ([]string, error) {
	var result []string
	for _, id := range ids {
		if err := validation.Uuid4(id); err == nil {
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

func BoolFromString(b string) (*upcloud.Boolean, error) {
	var result upcloud.Boolean
	if b == "true" {
		result = upcloud.FromBool(true)
	} else if b == "false" {
		result = upcloud.FromBool(false)
	} else {
		return nil, fmt.Errorf("invalid boolean value %s", b)
	}
	return &result, nil
}
