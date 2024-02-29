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
	labelSlice := make([]upcloud.Label, 0)

	for _, l := range in {
		label, err := stringToLabel(l)
		if err != nil {
			return nil, err
		}
		labelSlice = append(labelSlice, label)
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
