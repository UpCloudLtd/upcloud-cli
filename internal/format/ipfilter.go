package format

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/jedib0t/go-pretty/v6/text"
)

// IPFilter returns formatted IP filter.
func IPFilter(val interface{}) (text.Colors, string, error) {
	addresses, ok := val.([]string)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse IP addresses from %T, expected []string", val)
	}

	allowAll := false
	var strs []string
	for _, ipa := range addresses {
		if ipa == "0.0.0.0/0" {
			allowAll = true
		}

		strs = append(strs, ui.DefaultAddressColours.Sprint(ipa))
	}

	if allowAll {
		return nil, "all", nil
	}

	if len(addresses) == 0 {
		return nil, text.FgHiBlack.Sprint("none"), nil
	}

	return nil, strings.Join(strs, ",\n"), nil
}
