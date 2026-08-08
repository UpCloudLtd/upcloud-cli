package account

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
)

// LoginCommand creates the "account login" command
func LoginCommand() commands.Command {
	return &loginCommand{
		BaseCommand: commands.New(
			"login",
			"Configure an authentication token to the system keyring or credentials file",
			"upctl account login --with-token",
			"upctl account login --with-token --save-to-file",
		),
	}
}

type loginCommand struct {
	*commands.BaseCommand

	withToken  config.OptionalBoolean
	saveToFile bool
}

// InitCommand implements Command.InitCommand
func (s *loginCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	config.AddToggleFlag(fs, &s.withToken, "with-token", false, "Read token from standard input.")
	fs.BoolVar(&s.saveToFile, "save-to-file", false,
		"Save token to credentials file (~/.config/upcloud/credentials) instead of keyring when keyring is unavailable")
	s.AddFlags(fs)

	// Require the with-token flag until we support using browser to authenticate.
	commands.Must(s.Cobra().MarkFlagRequired("with-token"))
}

// DoesNotUseServices implements commands.OfflineCommand as this command does not use services
func (s *loginCommand) DoesNotUseServices() {}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *loginCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if s.withToken.Value() {
		return s.executeWithToken(exec)
	}

	return output.None{}, nil
}

func (s *loginCommand) executeWithToken(exec commands.Executor) (output.Output, error) {
	in := bufio.NewReader(s.Cobra().InOrStdin())
	token, err := in.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read token from standard input: %w", err)
	}
	token = strings.TrimSpace(token)

	// Try keyring first (unless explicitly using file)
	if !s.saveToFile {
		msg := "Saving provided token to the system keyring."
		exec.PushProgressStarted(msg)
		err = config.SaveTokenToKeyring(token)

		if err == nil {
			exec.PushProgressSuccess(msg)
			return output.None{}, nil
		}

		// Keyring failed
		if config.IsKeyringError(err) {
			var errMsg strings.Builder
			errMsg.WriteString("System keyring is not accessible.\n\n")

			// Provide system-specific hints
			if os.Getenv("WSL_DISTRO_NAME") != "" {
				errMsg.WriteString("  (WSL detected - keyring typically doesn't work in WSL)\n")
			} else if os.Getenv("SSH_CONNECTION") != "" {
				errMsg.WriteString("  (SSH session detected - keyring may not be available)\n")
			} else if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
				errMsg.WriteString("  (No display detected - keyring requires GUI or D-Bus)\n")
			}

			errMsg.WriteString("\nTo save your token, you can:\n")
			errMsg.WriteString("  1. Retry with file storage:\n")
			errMsg.WriteString("     upctl account login --with-token --save-to-file\n\n")
			errMsg.WriteString("  2. Set environment variable:\n")
			errMsg.WriteString("     export UPCLOUD_TOKEN=" + token + "\n\n")
			errMsg.WriteString("  3. Manually add to config file (~/.config/upctl.yaml):\n")
			errMsg.WriteString("     token: " + token + "\n\n")
			errMsg.WriteString("For security, tokens are not automatically saved to files without --save-to-file flag.")

			return commands.HandleError(exec, msg, fmt.Errorf("%s\n\nOriginal error: %w", errMsg.String(), err))
		}

		// Non-keyring error
		return commands.HandleError(exec, msg, err)
	}

	// User explicitly requested file storage
	msg := "Saving token to credentials file."
	exec.PushProgressStarted(msg)

	credPath, saveErr := config.SaveTokenToCredentialsFile(token)
	if saveErr != nil {
		errMsg := fmt.Sprintf("Failed to save token to file: %v\n\nYou can still use environment variable:\n  export UPCLOUD_TOKEN=%s", saveErr, token)
		return commands.HandleError(exec, msg, fmt.Errorf(errMsg))
	}

	exec.PushProgressSuccess(msg)
	fmt.Fprintf(s.Cobra().OutOrStderr(), "\nToken saved to: %s\n", credPath)
	fmt.Fprintf(s.Cobra().OutOrStderr(), "File permissions: 0600 (read/write for owner only)\n\n")
	fmt.Fprintf(s.Cobra().OutOrStderr(), "This credentials file can be shared by other UpCloud tools.\n")
	fmt.Fprintf(s.Cobra().OutOrStderr(), "You can now use: upctl account show\n")

	return output.None{}, nil
}
