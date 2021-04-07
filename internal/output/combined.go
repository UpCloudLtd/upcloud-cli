package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
)

// CombinedSection represents a single section of a combined output
type CombinedSection struct {
	Key      string
	Title    string
	Contents Output
}

// Combined represents multiple outputs combined and displayed sequentially (or wrapped into the same object)
type Combined []*CombinedSection

// MarshalJSON implements json.Marshaler
func (m Combined) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(flattenSections(m), "", "  ")
}

// MarshalYAML returns table output marshaled to YAML.
func (m Combined) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(flattenSections(m))
}

func flattenSections(m Combined) map[string]interface{} {
	out := map[string]interface{}{}
	for _, sec := range m {
		if sec == nil {
			continue
		}

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
	for _, sec := range m {
		if sec == nil {
			continue
		}

		marshaled, err := sec.Contents.MarshalHuman()
		if err != nil {
			return nil, err
		}
		if _, ok := sec.Contents.(Details); !ok && sec.Title != "" {
			// skip drawing title for details
			out = append(out, []byte(fmt.Sprintf("  %v\n", sec.Title))...)
		}
		out = append(out, marshaled...)
		out = append(out, []byte("\n\n")...)
	}
	return out, nil
}

// MarshalRawMap implements output.Output
func (m Combined) MarshalRawMap() (map[string]interface{}, error) {
	return nil, errors.New("multiple output cannot nest, raw output is undefined")
}
