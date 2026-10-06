package labels

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

func StringsToUpCloudLabelSlice(in []string) (*upcloud.LabelSlice, error) {
	upCloudlabelSlice := upcloud.LabelSlice{}

	for _, l := range in {
		label, err := stringToLabel(l)
		if err != nil {
			return nil, err
		}
		upCloudlabelSlice = append(upCloudlabelSlice, label)
	}

	return &upCloudlabelSlice, nil
}

func StringsToSliceOfLabels(in []string) ([]upcloud.Label, error) {
	return StringsToLabels(in, func(key, value string) upcloud.Label { return upcloud.Label{Key: key, Value: value} })
}

// StringsToLabels parses key=value strings into labels of any SDK type built by newLabel.
func StringsToLabels[T any](in []string, newLabel func(key, value string) T) ([]T, error) {
	labelSlice := make([]T, 0, len(in))

	for _, l := range in {
		label, err := stringToLabel(l)
		if err != nil {
			return nil, err
		}
		labelSlice = append(labelSlice, newLabel(label.Key, label.Value))
	}

	return labelSlice, nil
}

func stringToLabel(in string) (upcloud.Label, error) {
	split := strings.SplitN(in, "=", 2)
	if len(split) == 1 {
		return upcloud.Label{
			Key: split[0],
		}, nil
	}

	if len(split) == 2 {
		return upcloud.Label{
			Key:   split[0],
			Value: split[1],
		}, nil
	}

	return upcloud.Label{}, fmt.Errorf("invalid label: %s", in)
}
