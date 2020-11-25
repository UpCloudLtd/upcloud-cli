package commands

import (
	"encoding/csv"
	"strings"
)

func Parse(in string) ([]string, error) {
	var result []string
	reader := csv.NewReader(strings.NewReader(in))
	reader.LazyQuotes = true
	args, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for _, arg := range args {
		result = append(result, strings.Split("--"+arg, "=")...)
	}
	return result, nil
}
