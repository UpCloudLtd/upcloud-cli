package output

import (
	"gopkg.in/yaml.v3"
)

// JSONToYAML converts JSON bytes into YAML bytes. This allows using key names from JSON field tags also for YAML output. This has some side-effects (e.g., timestamps will be double-quoted in output) but this is lesser evil than adding yaml field tags everywhere in our Go types.
func JSONToYAML(jsonIn []byte) ([]byte, error) {
	if len(jsonIn) == 0 {
		return []byte{}, nil
	}

	var yamlObj interface{}
	if err := yaml.Unmarshal(jsonIn, &yamlObj); err != nil {
		return nil, err
	}
	return yaml.Marshal(yamlObj)
}
