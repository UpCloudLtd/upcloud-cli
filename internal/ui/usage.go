package ui

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var wrappingLineLength = 100

var templateFuncs = template.FuncMap{
	"trimTrailingWhitespaces": trimRightSpace,
	"rpad":                    rpad,
	"renderFlags":             formatFlags,
}

func formatFlags(fs *pflag.FlagSet) string {
	t := NewDataTable("flag", "usage")
	t.setStyle(styleFlagsTable())
	t.SetHeader(nil)
	fs.VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "help" {
			return
		}
		var flagText, flagUsage strings.Builder
		if flag.Shorthand != "" {
			flagText.WriteString(fmt.Sprintf("-%s, ", flag.Shorthand))
		}
		flagText.WriteString(fmt.Sprintf("--%s %s", flag.Name, flag.Value.Type()))
		flagUsage.WriteString(text.WrapSoft(flag.Usage, wrappingLineLength))
		def := flag.DefValue
		if strings.HasSuffix(flag.Value.Type(), "Slice") || strings.HasSuffix(flag.Value.Type(), "Array") {
			def = strings.TrimPrefix(def, "[")
			def = strings.TrimSuffix(def, "]")
		}
		if def != "" {
			flagUsage.WriteString(fmt.Sprintf("\nDefault: %s", def))
		}
		t.Append(table.Row{flagText.String(), flagUsage.String()})
	})
	return t.Render()
}

// Taken from cobra as they are private
func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

func rpad(s string, padding int) string {
	padTemplate := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(padTemplate, s)
}

// CommandUsageTemplate returns the template for usage
func CommandUsageTemplate() string {
	return `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
{{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Options:
{{renderFlags .LocalFlags}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Options:
{{renderFlags .InheritedFlags}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

// UsageFunc is used to override cobra's default usage func to get some more stylistic control
func UsageFunc(cmd *cobra.Command) error {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(cmd.UsageTemplate()))
	err := t.Execute(cmd.OutOrStdout(), cmd)
	return err
}
