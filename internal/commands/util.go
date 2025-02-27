package commands

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/validation"
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/crypto/ssh"
)

const (
	// FlagAnnotationNoFileCompletions is the annotation name to use for our flags that have no filename completions.
	FlagAnnotationNoFileCompletions = "upctl_flag_no_file_completions"
	// FlagAnnotationFixedCompletions is the annotation name to use for our flags that have a fixed set of completions.
	FlagAnnotationFixedCompletions = "upctl_flag_fixed_completions"
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

// SetDeprecationHelp hides a specific alias in the help output and prints a deprecation warning when used.
// Only works for primary commands, not subcommands.
func SetDeprecationHelp(cmd *cobra.Command, deprecatedAliases []string) {
	// Construct new alias list, excluding the deprecated aliases
	var filteredAliases []string
	for _, alias := range cmd.Aliases {
		if !slices.Contains(deprecatedAliases, alias) { // ✅ Using slices.Contains
			filteredAliases = append(filteredAliases, alias)
		}
	}

	// Update the alias list in the usage template **before** help is triggered
	originalTemplate := cmd.UsageTemplate()
	var modifiedTemplate string

	if len(filteredAliases) > 0 {
		modifiedTemplate = strings.ReplaceAll(originalTemplate, "{{.NameAndAliases}}{{end}}", "{{.Use}}, "+strings.Join(filteredAliases, ", ")+"{{end}}")
	} else {
		modifiedTemplate = strings.ReplaceAll(originalTemplate, "{{.NameAndAliases}}{{end}}", "{{.Use}}{{end}}")
	}

	// Apply the updated usage template **before help is called**
	cmd.SetUsageTemplate(modifiedTemplate)

	// Custom Help Function to Show Deprecation Warning
	cmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		if slices.Contains(deprecatedAliases, cmd.CalledAs()) {
			PrintDeprecationWarning(cmd.CalledAs(), cmd.Use)
		}

		// Print the help output with the modified alias list
		fmt.Println(cmd.UsageString())
	})

	// Intercept help execution using PersistentPreRunE
	originalPreRunE := cmd.PersistentPreRunE

	// Show deprecation message when upctl <command> help is called
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Call the parent command's PersistentPreRunE first
		if cmd.Parent() != nil && cmd.Parent().PersistentPreRunE != nil {
			if err := cmd.Parent().PersistentPreRunE(cmd.Parent(), args); err != nil {
				return err
			}
		}

		// Show deprecation warning if the alias was used
		if slices.Contains(deprecatedAliases, cmd.CalledAs()) {
			PrintDeprecationWarning(cmd.CalledAs(), cmd.Use)
		}

		// Call the original PersistentPreRunE (if defined)
		if originalPreRunE != nil {
			return originalPreRunE(cmd, args)
		}

		return nil
	}
}

// SetSubcommandDeprecationHelp detects the correct interface implementation and wraps the relevant execution function.
func SetSubcommandDeprecationHelp(cmd Command, aliases []string) {
	// Set a custom help function to display the warning for `-h` or `--help`
	cobraCmd := cmd.Cobra()
	originalHelpFunc := cobraCmd.HelpFunc()
	cobraCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		for _, alias := range aliases {
			if IsDeprecatedAliasUsed(alias) {
				PrintDeprecationWarning(alias, cmd.Parent().Use)
				break
			}
		}
		originalHelpFunc(cmd, args)
	})
}

// isDeprecatedAliasUsed checks if the deprecated alias was used in the command invocation
func IsDeprecatedAliasUsed(deprecatedAlias string) bool {
	if len(os.Args) < 2 {
		return false
	}

	for _, arg := range os.Args[1:] {
		if arg == deprecatedAlias {
			return true
		}
	}
	return false
}

// PrintDeprecationWarning prints a deprecation message
func PrintDeprecationWarning(deprecatedAlias, newCommand string) {
	fmt.Fprintf(os.Stderr, "Deprecation Warning: The alias '%s' is deprecated and will be removed in a future release.\n", deprecatedAlias)
	fmt.Fprintf(os.Stderr, "Please use '%s' instead.\n", newCommand)
}

func SetSubcommandExecutionDeprecationMessage(cmd Command, deprecatedParentAliases []string, mainParentAlias string) {
	parentCmd := cmd.Cobra().Parent()
	if parentCmd != nil {
		for _, deprecatedParentAlias := range deprecatedParentAliases {
			// Check if the parent was called using the deprecated alias
			if IsDeprecatedAliasUsed(deprecatedParentAlias) {
				PrintDeprecationWarning(deprecatedParentAlias, mainParentAlias)
				break
			}
		}
	}
}

// Must panics if the error is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
