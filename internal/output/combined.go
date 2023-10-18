package output

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
)

// CombinedSection represents a single section of a combined output
type CombinedSection struct {
	Key      string
	Title    string
	Contents Output
}

// Combined represents multiple outputs combined and displayed sequentially (or wrapped into the same object)
type Combined []CombinedSection

// MarshalJSON implements json.Marshaler
func (m Combined) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(flattenSections(m), "", "  ")
}

func flattenSections(m Combined) map[string]interface{} {
	out := map[string]interface{}{}

	for _, sec := range m {
		if sec.Key != "" {
			if _, ok := out[sec.Key]; ok {
				panic(fmt.Sprintf("duplicate section key '%v' in output", sec.Key))
			}
			if rowOut, err := sec.Contents.MarshalRawMap(); err != nil {
				panic(fmt.Sprintf("cannot marshal '%v' to raw output", sec.Key))
			} else {
				if len(rowOut) == 1 { // unwrap single ones - related to table; TODO: fix
					for _, v := range rowOut {
						out[sec.Key] = v
					}
				} else {
					out[sec.Key] = rowOut
				}
			}
		} else {
			if rowOut, err := sec.Contents.MarshalRawMap(); err != nil {
				panic(fmt.Sprintf("cannot marshal '%v' to raw output", sec.Key))
			} else {
				for k, v := range rowOut {
					out[k] = v
				}
			}
		}
	}

	return out
}

// MarshalHuman returns output in a human-readable form
func (m Combined) MarshalHuman() ([]byte, error) {
	out := []byte{}

	for i, sec := range m {
		marshaled, err := sec.Contents.MarshalHuman()
		if err != nil {
			return nil, err
		}
		if _, ok := sec.Contents.(Details); !ok && sec.Title != "" {
			// skip drawing title for details, as details handles its own title drawing
			// TODO: a bit confusing.. probably should refactor?
			out = append(out, []byte(fmt.Sprintf("  %v\n", ui.DefaultHeaderColours.Sprint(sec.Title)))...)
		}
		if _, ok := sec.Contents.(Table); ok {
			// this is a table, indent it
			marshaled = prefixLines(marshaled, "    ")
		}
		out = append(out, marshaled...)
		if i < len(m)-1 {
			// dont add newline after the last section
			out = append(out, []byte("\n")...)
		}
	}

	return out, nil
}

func prefixLines(marshaled []byte, s string) (out []byte) {
	padding := []byte(s)
	for _, b := range marshaled {
		if b == '\n' {
			out = append(append(out, b), padding...)
		} else {
			out = append(out, b)
		}
	}
	return
}

// MarshalRawMap implements output.Output
func (m Combined) MarshalRawMap() (map[string]interface{}, error) {
	return nil, errors.New("multiple output cannot nest, raw output is undefined")
}
