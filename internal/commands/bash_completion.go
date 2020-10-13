package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Override Cobra's custom completion function
// This fixes quoting of items with spaces
const goCustomCompletion = `__%[1]s_handle_go_custom_completion()
{
    __%[1]s_debug "${FUNCNAME[0]}: cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}"

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling ${words[0]} instead of directly %[1]s allows to handle aliases
    args=("${words[@]:1}")
    requestComp="${words[0]} __completeNoDesc ${args[@]:0:$((${#args[@]}-1))} $'$cur'"

    lastParam=${words[$((${#words[@]}-1))]}
    lastChar=${lastParam:$((${#lastParam}-1)):1}
    __%[1]s_debug "${FUNCNAME[0]}: lastParam ${lastParam}, lastChar ${lastChar}"

    if [ -z "${cur}" ] && [ "${lastChar}" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __%[1]s_debug "${FUNCNAME[0]}: Adding extra empty parameter"
        requestComp="${requestComp} \"\""
    fi

    __%[1]s_debug "${FUNCNAME[0]}: calling ${requestComp}"
    # Use eval to handle any environment variables and such
    out=$(eval "${requestComp}" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=${out##*:}
    # Remove the directive
    out=${out%%:*}
    if [ "${directive}" = "${out}" ]; then
        # There is not directive specified
        directive=0
    fi
    __%[1]s_debug "${FUNCNAME[0]}: the completion directive is: ${directive}"
    __%[1]s_debug "${FUNCNAME[0]}: the completions are: ${out[*]}"

    if [ $((directive & %[3]d)) -ne 0 ]; then
        # Error code.  No completion.
        __%[1]s_debug "${FUNCNAME[0]}: received error from custom completion go code"
        return
    else
        if [ $((directive & %[4]d)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __%[1]s_debug "${FUNCNAME[0]}: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & %[5]d)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __%[1]s_debug "${FUNCNAME[0]}: activating no file completion"
                compopt +o default
            fi
        fi

        local IFS=$'\n'
        COMPREPLY=($out)
    fi
}`

func CustomBashCompletionFunc(name string) string {
	return fmt.Sprintf(goCustomCompletion, name,
		cobra.ShellCompNoDescRequestCmd,
		cobra.ShellCompDirectiveError,
		cobra.ShellCompDirectiveNoSpace,
		cobra.ShellCompDirectiveNoFileComp)
}
