package cmd

import (
	"io/ioutil"
	"os"
	"strings"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func readFileLines(filename string) ([]string, error) {
	content, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	rawLines := strings.Split(content, "\n")
	lines := []string{}
	for _, line := range rawLines {
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func maxInt(numbers ...int) int {
	max := numbers[0]
	for _, number := range numbers {
		if number > max {
			max = number
		}
	}

	return max
}
