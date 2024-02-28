package commands

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/validation"

	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/crypto/ssh"
)

// ParseN parses a complex, querystring-type argument from `in` and splits values to `n` amount of substrings
// e.g. with `n` 2: `--foo bar=baz,flop=flip=1` returns `[]string{"bar","baz","flop","flip=1"}`
func ParseN(in string, n int) ([]string, error) {
	var result []string
	reader := csv.NewReader(strings.NewReader(in))
	args, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for _, arg := range args {
		result = append(result, strings.SplitN("--"+arg, "=", n)...)
	}
	return result, nil
}

// Parse calls `ParseN()` with `n` -1:
// eg. `--foo bar=baz,flop=flip` returns `[]string{"bar","baz","flop","flip"}` and
// `--foo bar=baz,flop=flip=1` returns `[]string{"bar","baz","flop","flip","1"}`
func Parse(in string) ([]string, error) {
	return ParseN(in, -1)
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

// WrapLongDescription wraps Long description messages at 80 characters and removes trailing whitespace from the message.
func WrapLongDescription(message string) string {
	re := regexp.MustCompile(` +\n`)
	wrapped := text.WrapSoft(message, 80)
	return re.ReplaceAllString(wrapped, "\n")
}

// ParseSSHKeys parses strings that can be either actual public keys
// or file names referring public key files.
func ParseSSHKeys(sshKeys []string) ([]string, error) {
	var allSSHKeys []string
	for _, keyOrFile := range sshKeys {
		if strings.HasPrefix(keyOrFile, "ssh-") {
			if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(keyOrFile)); err != nil {
				return nil, fmt.Errorf("invalid ssh key %q: %v", keyOrFile, err)
			}
			allSSHKeys = append(allSSHKeys, keyOrFile)
			continue
		}

		f, err := os.Open(keyOrFile)
		if err != nil {
			return nil, err
		}

		rdr := bufio.NewScanner(f)
		for rdr.Scan() {
			if _, _, _, _, err := ssh.ParseAuthorizedKey(rdr.Bytes()); err != nil {
				_ = f.Close()
				return nil, fmt.Errorf("invalid ssh key %q in file %s: %v", rdr.Text(), keyOrFile, err)
			}
			allSSHKeys = append(allSSHKeys, rdr.Text())
		}
		_ = f.Close()
	}

	return allSSHKeys, nil
}
