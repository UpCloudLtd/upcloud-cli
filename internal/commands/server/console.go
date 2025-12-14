package server

import (
	"fmt"
	"os"
	osExec "os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type consoleCommand struct {
	*commands.BaseCommand
	viewer       string
	fullscreen   bool
	viewOnly     bool
	showPassword bool
	completion.Server
	resolver.CachingServer
}

// ConsoleCommand creates the "server console" command
func ConsoleCommand() commands.Command {
	return &consoleCommand{
		BaseCommand: commands.New(
			"console",
			"Connect to server VNC console",
			"upctl server console 00038afc-e100-4e91-9d28-b9c463e7e9b4",
			"upctl server console myserver --fullscreen",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *consoleCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.viewer, "viewer", "", "VNC client to use (tigervnc, realvnc, remmina, macos)")
	flags.BoolVar(&s.fullscreen, "fullscreen", false, "Start in fullscreen mode (TigerVNC, Remmina)")
	flags.BoolVar(&s.viewOnly, "view-only", false, "View only, no input (TigerVNC, RealVNC)")
	flags.BoolVar(&s.showPassword, "show-password", false, "Display VNC password (required for Remmina and macOS Screen Sharing)")
	s.AddFlags(flags)

	commands.Must(s.Cobra().RegisterFlagCompletionFunc("viewer", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"tigervnc", "realvnc", "remmina", "macos"}, cobra.ShellCompDirectiveNoFileComp
	}))
}

// Execute implements commands.MultipleArgumentCommand
func (s *consoleCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	msg := fmt.Sprintf("Connecting to VNC console on server %v", uuid)
	exec.PushProgressStarted(msg)

	// Get server details
	svc := exec.All()
	serverDetails, err := svc.GetServerDetails(exec.Context(), &request.GetServerDetailsRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	// Check if VNC is enabled
	if !serverDetails.RemoteAccessEnabled.Bool() {
		return nil, fmt.Errorf("VNC remote access is not enabled on server %s. Enable it with 'upctl server modify %s --remote-access-enabled yes --remote-access-type vnc'", uuid, uuid)
	}

	if serverDetails.RemoteAccessType != upcloud.RemoteAccessTypeVNC {
		return nil, fmt.Errorf("server is configured for %s, not VNC", serverDetails.RemoteAccessType)
	}

	// Extract VNC credentials
	host := serverDetails.RemoteAccessHost
	port := serverDetails.RemoteAccessPort
	password := serverDetails.RemoteAccessPassword

	if host == "" || port == 0 {
		return nil, fmt.Errorf("VNC connection details not available")
	}

	// Detect or use specified VNC client
	client, err := s.detectVNCClient()
	if err != nil {
		return nil, err
	}

	// Display password if explicitly requested or if using a client that requires it
	if s.showPassword {
		exec.PushProgressSuccess(fmt.Sprintf("Launching %s to connect to %s...\nVNC Password: %s", client.name, serverDetails.Title, password))
	} else if client.name == "remmina" || client.name == "macos" {
		exec.PushProgressSuccess(fmt.Sprintf("Launching %s to connect to %s...\nNote: Use --show-password flag to display the VNC password for manual entry.", client.name, serverDetails.Title))
	} else {
		exec.PushProgressSuccess(fmt.Sprintf("Launching %s to connect to %s...", client.name, serverDetails.Title))
	}

	// Create secure temporary password file
	passFile, cleanup, err := createSecurePasswordFile(password)
	if err != nil {
		return nil, fmt.Errorf("failed to create password file: %w", err)
	}
	defer cleanup()

	// Build and execute VNC client command
	args := client.buildArgs(host, port, passFile, s.fullscreen, s.viewOnly)
	vncCmd := osExec.Command(client.executable, args...)
	vncCmd.Stdout = os.Stdout
	vncCmd.Stderr = os.Stderr

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		if vncCmd.Process != nil {
			vncCmd.Process.Signal(os.Interrupt)
		}
		cleanup()
	}()

	// Run VNC client
	if err := vncCmd.Run(); err != nil {
		return nil, fmt.Errorf("VNC client failed: %w", err)
	}

	return output.None{}, nil
}

// vncClient represents a VNC client configuration
type vncClient struct {
	name       string
	executable string
	buildArgs  func(host string, port int, passFile string, fullscreen, viewOnly bool) []string
}

// detectVNCClient finds the first available VNC client
func (s *consoleCommand) detectVNCClient() (*vncClient, error) {
	clients := s.getAvailableVNCClients()

	// If user specified a viewer, try to find it
	if s.viewer != "" {
		for _, client := range clients {
			if client.name == s.viewer {
				if _, err := osExec.LookPath(client.executable); err == nil {
					return &client, nil
				}
				return nil, fmt.Errorf("VNC client '%s' not found in PATH", s.viewer)
			}
		}
		return nil, fmt.Errorf("unknown VNC client '%s'", s.viewer)
	}

	// Auto-detect first available client
	for _, client := range clients {
		if _, err := osExec.LookPath(client.executable); err == nil {
			return &client, nil
		}
	}

	// No client found - provide platform-specific installation instructions
	return nil, fmt.Errorf("no VNC client found. Install one:\n%s", getInstallInstructions())
}

// getInstallInstructions returns platform-specific VNC client installation instructions
func getInstallInstructions() string {
	switch runtime.GOOS {
	case "darwin":
		return "  macOS:\n" +
			"    Built-in Screen Sharing is available (no installation needed)\n" +
			"    Or install TigerVNC: brew install --cask tiger-vnc-viewer"
	case "linux":
		// Check for specific distributions
		distro := detectLinuxDistro()
		switch distro {
		case "ubuntu", "debian":
			return "  Ubuntu/Debian:\n" +
				"    sudo apt-get install -y remmina-plugin-vnc tigervnc-viewer"
		case "fedora", "rhel", "centos":
			return "  Fedora/RHEL/CentOS:\n" +
				"    sudo dnf install -y remmina tigervnc"
		case "arch":
			return "  Arch Linux:\n" +
				"    sudo pacman -S remmina tigervnc"
		default:
			return "  Linux (generic):\n" +
				"    Debian/Ubuntu:     sudo apt-get install -y remmina-plugin-vnc tigervnc-viewer\n" +
				"    Fedora/RHEL:       sudo dnf install -y remmina tigervnc\n" +
				"    Arch:              sudo pacman -S remmina tigervnc\n" +
				"    Or use your package manager to install remmina or tigervnc-viewer"
		}
	case "windows":
		return "  Windows:\n" +
			"    Download TigerVNC from: https://github.com/TigerVNC/tigervnc/releases\n" +
			"    Or RealVNC from: https://www.realvnc.com/en/connect/download/viewer/"
	default:
		return "  Install TigerVNC or RealVNC for your platform"
	}
}

// detectLinuxDistro attempts to detect the Linux distribution
func detectLinuxDistro() string {
	// Try reading /etc/os-release
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}

	content := string(data)
	if strings.Contains(strings.ToLower(content), "ubuntu") {
		return "ubuntu"
	}
	if strings.Contains(strings.ToLower(content), "debian") {
		return "debian"
	}
	if strings.Contains(strings.ToLower(content), "fedora") {
		return "fedora"
	}
	if strings.Contains(strings.ToLower(content), "rhel") || strings.Contains(strings.ToLower(content), "red hat") {
		return "rhel"
	}
	if strings.Contains(strings.ToLower(content), "centos") {
		return "centos"
	}
	if strings.Contains(strings.ToLower(content), "arch") {
		return "arch"
	}

	return ""
}

// getAvailableVNCClients returns VNC clients in platform-specific preferred order
func (s *consoleCommand) getAvailableVNCClients() []vncClient {
	var clients []vncClient

	// Platform-specific ordering: prefer native/built-in clients first
	switch runtime.GOOS {
	case "darwin":
		// macOS: prefer built-in Screen Sharing, then TigerVNC
		clients = append(clients, vncClient{
			name:       "macos",
			executable: "open",
			buildArgs: func(host string, port int, passFile string, fullscreen, viewOnly bool) []string {
				// macOS Screen Sharing prompts for password interactively
				// Note: fullscreen and viewOnly are not supported by macOS Screen Sharing
				return []string{fmt.Sprintf("vnc://%s:%d", host, port)}
			},
		})
		clients = append(clients, tigervncClient())
		clients = append(clients, realvncClient())

	case "linux":
		// Linux: prefer Remmina (common on GNOME/Ubuntu), then TigerVNC
		clients = append(clients, vncClient{
			name:       "remmina",
			executable: "remmina",
			buildArgs: func(host string, port int, passFile string, fullscreen, viewOnly bool) []string {
				// Remmina will prompt for password interactively
				args := []string{}
				if fullscreen {
					args = append(args, "--enable-fullscreen")
				}
				args = append(args, "-c", fmt.Sprintf("vnc://%s:%d", host, port))
				// Note: viewOnly flag is not supported by Remmina (no command-line option available)
				return args
			},
		})
		clients = append(clients, tigervncClient())
		clients = append(clients, realvncClient())

	default:
		// Windows and others: TigerVNC first, then RealVNC
		clients = append(clients, tigervncClient())
		clients = append(clients, realvncClient())
	}

	return clients
}

// tigervncClient returns TigerVNC client configuration
func tigervncClient() vncClient {
	return vncClient{
		name:       "tigervnc",
		executable: "vncviewer",
		buildArgs: func(host string, port int, passFile string, fullscreen, viewOnly bool) []string {
			args := []string{}
			if fullscreen {
				args = append(args, "-fullscreen")
			}
			if viewOnly {
				args = append(args, "-viewonly")
			}
			args = append(args, "-passwd", passFile, fmt.Sprintf("%s::%d", host, port))
			return args
		},
	}
}

// realvncClient returns RealVNC client configuration
func realvncClient() vncClient {
	return vncClient{
		name:       "realvnc",
		executable: "vncviewer",
		buildArgs: func(host string, port int, passFile string, fullscreen, viewOnly bool) []string {
			args := []string{"-passwd", passFile}
			if viewOnly {
				args = append(args, "-viewonly")
			}
			args = append(args, fmt.Sprintf("%s::%d", host, port))
			return args
		},
	}
}

// createSecurePasswordFile creates a temporary file with VNC password
func createSecurePasswordFile(password string) (string, func(), error) {
	// Create temp directory with restrictive permissions
	tmpDir, err := os.MkdirTemp("", "upctl-vnc-*")
	if err != nil {
		return "", nil, err
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	// Set directory permissions to 0700 (owner only)
	if err := os.Chmod(tmpDir, 0700); err != nil {
		cleanup()
		return "", nil, err
	}

	// Create password file
	passFile := filepath.Join(tmpDir, "password")

	// Write password as plain text (TigerVNC -passwd accepts plain text)
	if err := os.WriteFile(passFile, []byte(password), 0600); err != nil {
		cleanup()
		return "", nil, err
	}

	return passFile, cleanup, nil
}
