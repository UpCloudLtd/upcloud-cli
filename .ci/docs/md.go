package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// formatFlags formats options and aliases as list of code elements. E.g., ["--help", "-h"] becomes "`--help`, `-h`".
func formatFlags(input []string) string {
	return fmt.Sprintf("`%s`", strings.Join(input, "`, `"))
}

// formatDescription replaces line breaks in descriptions with <br/> elements to avoid breaking table formatting.
func formatDescription(input string) string {
	return strings.Replace(input, "\n", "<br/>", -1)
}

// sprintOptionsTable renders options in table to make them more accessible.
func sprintOptionsTable(flags *pflag.FlagSet) string {
	output := "| Option | Description |\n"
	output += "| ------ | ----------- |\n"

	flags.VisitAll(func(flag *pflag.Flag) {
		flags := []string{"--" + flag.Name}
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
			flags = append(flags, "-"+flag.Shorthand)
		}

		description := flag.Usage
		if flag.DefValue != "" {
			description += fmt.Sprintf("\nDefault: `%s`", flag.DefValue)
		}

		row := fmt.Sprintf("| %s | %s |\n", formatFlags(flags), formatDescription(description))
		output += row
	})

	return output
}

func sprintOptions(cmd *cobra.Command, name string) string {
	output := ""

	flags := cmd.NonInheritedFlags()
	if flags.HasAvailableFlags() {
		output += "## Options\n\n"
		output += sprintOptionsTable(flags)
		output += "\n"
	}

	parentFlags := cmd.InheritedFlags()
	if parentFlags.HasAvailableFlags() {
		output += "## Global options\n\n"
		output += sprintOptionsTable(parentFlags)
		output += "\n"
	}

	return output
}

func isHiddenCommand(cmd *cobra.Command) bool {
	return !cmd.IsAvailableCommand() || cmd.IsAdditionalHelpTopicCommand()
}

func hasChildCommands(cmd *cobra.Command) bool {
	for _, c := range cmd.Commands() {
		if isHiddenCommand(c) {
			continue
		}
		return true
	}
	return false
}

func hasRelatedCommands(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}
	return hasChildCommands(cmd)
}

func getFilename(cmd *cobra.Command) string {
	cmdpath := cmd.CommandPath()
	if hasChildCommands(cmd) && strings.Contains(cmdpath, " ") {
		cmdpath = cmdpath + " index"
	}
	basename := strings.Replace(cmdpath, " ", "_", 1)
	return strings.ReplaceAll(basename, " ", "/") + ".md"
}

func getRelativeLink(current, target *cobra.Command) string {
	depth := strings.Count(getFilename(current), "/")
	return strings.Repeat("../", depth) + getFilename(target)
}

func genMarkdown(cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("# " + name + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString(cmd.Long + "\n\n")
	} else {
		buf.WriteString(cmd.Short + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.UseLine()))
	}

	ctx := cmd.Context()
	if ctx != nil {
		if wrapper := ctx.Value("command"); wrapper != nil {
			if c, ok := wrapper.(*commands.BaseCommand); ok {
				if aliases := c.Aliases(); len(aliases) > 0 {
					buf.WriteString("## Aliases\n\n")
					buf.WriteString(fmt.Sprintf("%s\n\n", formatFlags(aliases)))
				}
			}
		}
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("## Examples\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	buf.WriteString(sprintOptions(cmd, name))

	if hasRelatedCommands(cmd) {
		buf.WriteString("## Related commands\n\n")
		buf.WriteString("| Command | Description |\n")
		buf.WriteString("| ------- | ----------- |\n")

		if cmd.HasParent() {
			parent := cmd.Parent()
			pname := parent.CommandPath()
			link := getRelativeLink(cmd, parent)
			buf.WriteString(fmt.Sprintf("| [%s](%s) | %s |\n", pname, link, formatDescription(parent.Short)))
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			link := getRelativeLink(cmd, child)
			buf.WriteString(fmt.Sprintf("| [%s](%s) | %s |\n", cname, link, formatDescription(child.Short)))
		}
		buf.WriteString("\n")
	}

	_, err := buf.WriteTo(w)
	return err
}

func genMarkdownTree(cmd *cobra.Command, dir string, filePrepender func(*cobra.Command) string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := genMarkdownTree(c, dir, filePrepender); err != nil {
			return err
		}
	}

	basename := getFilename(cmd)
	filename := filepath.Join(dir, basename)
	directory := path.Dir(filename)
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.WriteString(f, filePrepender(cmd)); err != nil {
		return err
	}
	if err := genMarkdown(cmd, f); err != nil {
		return err
	}
	return nil
}
