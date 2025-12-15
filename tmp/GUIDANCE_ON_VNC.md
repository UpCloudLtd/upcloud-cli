# Guidance on VNC Console Access for UpCloud CLI

**Purpose:** Enable secure VNC console access to UpCloud servers via the `upctl` CLI without exposing credentials on command line

**Author:** Assistant
**Date:** 2025-12-13
**Status:** Implementation Guide

---

## Executive Summary

UpCloud provides VNC console access to servers through their API. The VNC credentials (host, port, password) are available in the server details. This guide shows how to implement secure VNC console access in `upctl` that:

1. ✅ **Never exposes passwords in command line arguments**
2. ✅ **Uses temporary files or environment variables for credentials**
3. ✅ **Supports multiple VNC clients across platforms**
4. ✅ **Cleans up credentials after connection**
5. ✅ **Provides a seamless user experience**

---

## UpCloud VNC API Support

### Already Implemented in upcloud-go-api

From `upcloud-go-api/v8/upcloud/server.go`:

```go
type ServerDetails struct {
    // ... other fields ...
    RemoteAccessEnabled  Boolean  `json:"remote_access_enabled"`
    RemoteAccessType     string   `json:"remote_access_type"`      // "vnc"
    RemoteAccessHost     string   `json:"remote_access_host"`      // e.g., "fi-hel1.vnc.upcloud.com"
    RemoteAccessPassword string   `json:"remote_access_password"`  // Generated password
    RemoteAccessPort     int      `json:"remote_access_port"`      // VNC port number
}
```

### Enabling VNC on Server Creation

```go
serverDetails, err := svc.CreateServer(ctx, &request.CreateServerRequest{
    // ... other fields ...
    RemoteAccessEnabled:  upcloud.True,
    RemoteAccessType:     upcloud.RemoteAccessTypeVNC,
    RemoteAccessPassword: "custom-password",  // Optional - auto-generated if omitted
})
```

### Retrieving VNC Credentials

```go
serverDetails, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{
    UUID: serverUUID,
})

if serverDetails.RemoteAccessEnabled.Bool() {
    host := serverDetails.RemoteAccessHost
    port := serverDetails.RemoteAccessPort
    password := serverDetails.RemoteAccessPassword
}
```

---

## VNC Clients: Cross-Platform Options

### TigerVNC (Recommended - Open Source)

**Platform Support:** Linux, Windows, macOS
**License:** GPL
**Installation:**

```bash
# Ubuntu/Debian
sudo apt install tigervnc-viewer

# macOS (Homebrew)
brew install --cask tiger-vnc-viewer

# Windows
# Download from https://github.com/TigerVNC/tigervnc/releases
```

**Command Line Usage:**

```bash
# Basic connection
vncviewer <host>:<display>

# With password from stdin (SECURE)
echo "$VNC_PASSWORD" | vncviewer -autopass <host>:<display>

# With password file (SECURE)
vncviewer -passwd /tmp/vnc-password-file <host>:<display>

# Full screen mode
vncviewer -fullscreen <host>:<display>
```

**Key Parameters:**

| Parameter | Description |
|-----------|-------------|
| `-autopass` | Read password from stdin (secure) |
| `-passwd <file>` | Read password from file (secure) |
| `-fullscreen` | Start in full screen mode |
| `-viewonly` | View only, no input |
| `-shared` | Allow multiple connections |
| `-depth <bits>` | Color depth (8, 16, 24) |
| `-quality <level>` | Compression quality (0-9) |

**Environment Variables:**

```bash
VNC_USERNAME  # VNC username (if required)
VNC_PASSWORD  # VNC password
```

### RealVNC Viewer

**Platform Support:** Linux, Windows, macOS, iOS, Android
**License:** Free tier available, commercial for advanced features
**Installation:**

```bash
# Download from https://www.realvnc.com/en/connect/download/viewer/
```

**Command Line Usage:**

```bash
# Basic connection
vncviewer <host>:<display>

# With parameter file (SECURE)
vncviewer -config /tmp/vnc-config <host>:<display>
```

### UltraVNC (Windows)

**Platform Support:** Windows
**License:** GPL

**Command Line Usage:**

```bash
# With password parameter (LESS SECURE - visible in process list)
vncviewer.exe <host>:<display> -password <password>

# Better: use password file
vncviewer.exe <host>:<display> -passwordfile vnc.pwd
```

### Built-in VNC Clients

**macOS:**
```bash
# Built-in Screen Sharing app
open vnc://<host>:<port>
# Note: Will prompt for password interactively
```

**Linux (GNOME):**
```bash
# Remmina (often pre-installed)
remmina -c vnc://<host>:<display>
```

---

## Security Best Practices

### ❌ NEVER DO THIS (Insecure)

```bash
# Password visible in command line arguments
vncviewer server.com:5901 -password "mypassword123"

# Password visible in process list
ps aux | grep vncviewer
# Shows: vncviewer server.com:5901 -password mypassword123
```

### ✅ DO THIS (Secure)

#### Method 1: Password from stdin with `-autopass`

```bash
# Password never appears in process list
echo "$VNC_PASSWORD" | vncviewer -autopass server.com:5901
```

#### Method 2: Temporary password file

```bash
# Create secure temporary file
PASSFILE=$(mktemp -t vnc-pass.XXXXXX)
chmod 600 "$PASSFILE"
echo "$VNC_PASSWORD" > "$PASSFILE"

# Use password file
vncviewer -passwd "$PASSFILE" server.com:5901

# Clean up
rm -f "$PASSFILE"
```

#### Method 3: Environment variable

```bash
# Set environment variable (only visible to current process)
export VNC_PASSWORD="secret123"
echo "$VNC_PASSWORD" | vncviewer -autopass server.com:5901
unset VNC_PASSWORD
```

---

## Recommended CLI Implementation

### Proposed Command Syntax

```bash
# Connect to VNC console
upctl server console <server-uuid>

# Options
upctl server console <server-uuid> --fullscreen
upctl server console <server-uuid> --viewer tigervnc
upctl server console <server-uuid> --viewer realvnc
upctl server console <server-uuid> --view-only
```

### Implementation Architecture

```
┌─────────────────┐
│  upctl console  │
└────────┬────────┘
         │
         ├─ 1. Get server details via API
         ├─ 2. Extract VNC credentials
         ├─ 3. Create secure temp password file
         ├─ 4. Detect/select VNC client
         ├─ 5. Launch VNC client with secure credentials
         └─ 6. Clean up temp files on exit
```

### Complete Implementation Example

```go
package console

import (
    "context"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "os/signal"
    "path/filepath"
    "runtime"
    "syscall"

    "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
    "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
    "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
)

// VNCClient represents a VNC client configuration
type VNCClient struct {
    Name       string
    Executable string
    Args       func(host string, port int, passFile string) []string
}

// GetAvailableVNCClients returns VNC clients available on the system
func GetAvailableVNCClients() []VNCClient {
    clients := []VNCClient{
        {
            Name:       "TigerVNC",
            Executable: "vncviewer",
            Args: func(host string, port int, passFile string) []string {
                return []string{"-passwd", passFile, fmt.Sprintf("%s:%d", host, port)}
            },
        },
        {
            Name:       "RealVNC",
            Executable: "vncviewer",
            Args: func(host string, port int, passFile string) []string {
                // RealVNC also supports -passwd
                return []string{"-passwd", passFile, fmt.Sprintf("%s:%d", host, port)}
            },
        },
    }

    // macOS-specific
    if runtime.GOOS == "darwin" {
        clients = append(clients, VNCClient{
            Name:       "macOS Screen Sharing",
            Executable: "open",
            Args: func(host string, port int, passFile string) []string {
                // Note: macOS open vnc:// will prompt for password
                return []string{fmt.Sprintf("vnc://%s:%d", host, port)}
            },
        })
    }

    return clients
}

// DetectVNCClient finds the first available VNC client
func DetectVNCClient() (*VNCClient, error) {
    clients := GetAvailableVNCClients()

    for _, client := range clients {
        if _, err := exec.LookPath(client.Executable); err == nil {
            return &client, nil
        }
    }

    return nil, fmt.Errorf("no VNC client found. Please install TigerVNC or RealVNC")
}

// ConnectVNC establishes a VNC connection to the server
func ConnectVNC(ctx context.Context, svc *service.Service, serverUUID string, options *VNCOptions) error {
    // 1. Get server details
    serverDetails, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{
        UUID: serverUUID,
    })
    if err != nil {
        return fmt.Errorf("failed to get server details: %w", err)
    }

    // 2. Check if VNC is enabled
    if !serverDetails.RemoteAccessEnabled.Bool() {
        return fmt.Errorf("VNC remote access is not enabled on this server. Enable it with 'upctl server modify %s --enable-vnc'", serverUUID)
    }

    if serverDetails.RemoteAccessType != upcloud.RemoteAccessTypeVNC {
        return fmt.Errorf("server is configured for %s, not VNC", serverDetails.RemoteAccessType)
    }

    // 3. Extract VNC credentials
    host := serverDetails.RemoteAccessHost
    port := serverDetails.RemoteAccessPort
    password := serverDetails.RemoteAccessPassword

    if host == "" || port == 0 {
        return fmt.Errorf("VNC connection details not available")
    }

    // 4. Create secure temporary password file
    passFile, err := createSecurePasswordFile(password)
    if err != nil {
        return fmt.Errorf("failed to create password file: %w", err)
    }
    defer os.Remove(passFile) // Clean up on exit

    // 5. Detect or use specified VNC client
    var client *VNCClient
    if options.Viewer != "" {
        client = findClientByName(options.Viewer)
        if client == nil {
            return fmt.Errorf("VNC client '%s' not found", options.Viewer)
        }
    } else {
        client, err = DetectVNCClient()
        if err != nil {
            return err
        }
    }

    fmt.Printf("Connecting to %s via VNC using %s...\n", serverDetails.Title, client.Name)
    fmt.Printf("VNC Server: %s:%d\n", host, port)

    // 6. Build command arguments
    args := client.Args(host, port, passFile)

    // Add fullscreen flag if requested
    if options.Fullscreen && client.Name == "TigerVNC" {
        args = append([]string{"-fullscreen"}, args...)
    }

    // 7. Launch VNC client
    cmd := exec.Command(client.Executable, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Handle Ctrl+C gracefully
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigChan
        if cmd.Process != nil {
            cmd.Process.Signal(os.Interrupt)
        }
        os.Remove(passFile) // Ensure cleanup on interrupt
    }()

    // 8. Run and wait
    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("VNC client failed: %w", err)
    }

    return nil
}

// createSecurePasswordFile creates a temporary file with VNC password
func createSecurePasswordFile(password string) (string, error) {
    // Create temp directory with restrictive permissions
    tmpDir, err := ioutil.TempDir("", "upctl-vnc-*")
    if err != nil {
        return "", err
    }

    // Set directory permissions to 0700 (owner only)
    if err := os.Chmod(tmpDir, 0700); err != nil {
        os.RemoveAll(tmpDir)
        return "", err
    }

    // Create password file
    passFile := filepath.Join(tmpDir, "password")

    // For TigerVNC, we need to create a VNC password file format
    // For simplicity, we'll write plain text and rely on file permissions
    // Note: TigerVNC's vncpasswd uses DES encryption, but plain text works with -passwd
    if err := ioutil.WriteFile(passFile, []byte(password), 0600); err != nil {
        os.RemoveAll(tmpDir)
        return "", err
    }

    return passFile, nil
}

// VNCOptions configures VNC connection behavior
type VNCOptions struct {
    Viewer     string // Specific VNC client to use
    Fullscreen bool   // Start in fullscreen mode
    ViewOnly   bool   // View only, no input
}

func findClientByName(name string) *VNCClient {
    clients := GetAvailableVNCClients()
    for _, client := range clients {
        if client.Name == name {
            return &client
        }
    }
    return nil
}
```

### Alternative: Using stdin with `-autopass`

For TigerVNC, we can avoid temporary files entirely:

```go
func ConnectVNCWithStdin(host string, port int, password string) error {
    cmd := exec.Command("vncviewer", "-autopass", fmt.Sprintf("%s:%d", host, port))

    // Create pipe for password
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return err
    }

    // Start the command
    if err := cmd.Start(); err != nil {
        return err
    }

    // Write password to stdin
    _, err = stdin.Write([]byte(password + "\n"))
    stdin.Close()
    if err != nil {
        return err
    }

    // Wait for completion
    return cmd.Wait()
}
```

---

## Command Integration

### File Structure

```
internal/commands/server/
├── console/
│   ├── console.go          # Main console command
│   ├── vnc.go              # VNC connection logic
│   ├── vnc_clients.go      # VNC client detection
│   └── vnc_test.go         # Tests
```

### Command Registration

In `internal/commands/server/server.go`:

```go
import (
    "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/server/console"
)

func BuildCommands() []commands.Command {
    return []commands.Command{
        // ... existing commands
        console.BaseConsoleCommand(),
    }
}
```

### Console Command Implementation

```go
package console

import (
    "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
    "github.com/spf13/pflag"
)

type consoleCommand struct {
    *commands.BaseCommand
    viewer     string
    fullscreen bool
    viewOnly   bool
}

func BaseConsoleCommand() commands.Command {
    return &consoleCommand{
        BaseCommand: commands.New(
            "console",
            "Connect to server VNC console",
            "upctl server console <uuid>",
            "upctl server console 00000000-0000-0000-0000-000000000000",
        ),
    }
}

func (c *consoleCommand) InitCommand() {
    flags := &pflag.FlagSet{}
    flags.StringVar(&c.viewer, "viewer", "", "VNC client to use (tigervnc, realvnc)")
    flags.BoolVar(&c.fullscreen, "fullscreen", false, "Start in fullscreen mode")
    flags.BoolVar(&c.viewOnly, "view-only", false, "View only, no input")
    c.AddFlags(flags)
}

func (c *consoleCommand) Execute(exec commands.Executor) error {
    svc := exec.All()

    serverUUID := c.Args()[0]

    options := &VNCOptions{
        Viewer:     c.viewer,
        Fullscreen: c.fullscreen,
        ViewOnly:   c.viewOnly,
    }

    return ConnectVNC(exec.Context(), svc, serverUUID, options)
}
```

---

## User Experience

### Successful Connection

```bash
$ upctl server console 00000000-0000-0000-0000-000000000000

Connecting to my-server via VNC using TigerVNC...
VNC Server: fi-hel1.vnc.upcloud.com:5901

[VNC window opens]
```

### Error Cases

```bash
# VNC not enabled
$ upctl server console 00000000-0000-0000-0000-000000000000
Error: VNC remote access is not enabled on this server.
Enable it with 'upctl server modify 00000000-0000-0000-0000-000000000000 --enable-vnc'

# No VNC client found
$ upctl server console 00000000-0000-0000-0000-000000000000
Error: no VNC client found. Please install TigerVNC or RealVNC

Install instructions:
  Ubuntu/Debian: sudo apt install tigervnc-viewer
  macOS:         brew install --cask tiger-vnc-viewer
  Windows:       Download from https://github.com/TigerVNC/tigervnc/releases
```

### Help Output

```bash
$ upctl server console --help

Connect to server VNC console

Usage:
  upctl server console <uuid> [flags]

Flags:
      --viewer string       VNC client to use (tigervnc, realvnc)
      --fullscreen         Start in fullscreen mode
      --view-only          View only, no input

Examples:
  upctl server console 00000000-0000-0000-0000-000000000000
  upctl server console 00000000-0000-0000-0000-000000000000 --fullscreen
  upctl server console 00000000-0000-0000-0000-000000000000 --viewer tigervnc
```

---

## Testing

### Unit Tests

```go
func TestDetectVNCClient(t *testing.T) {
    client, err := DetectVNCClient()

    // Should find at least one client on developer machine
    if err != nil {
        t.Skip("No VNC client installed on test machine")
    }

    assert.NotNil(t, client)
    assert.NotEmpty(t, client.Name)
}

func TestCreateSecurePasswordFile(t *testing.T) {
    password := "test-password-123"

    passFile, err := createSecurePasswordFile(password)
    assert.NoError(t, err)
    defer os.Remove(passFile)

    // Check file permissions
    info, err := os.Stat(passFile)
    assert.NoError(t, err)
    assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

    // Verify content
    content, err := ioutil.ReadFile(passFile)
    assert.NoError(t, err)
    assert.Equal(t, password, string(content))
}
```

### Integration Tests

Test against actual VNC server (requires credentials):

```go
func TestVNCConnection(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Requires UPCLOUD_GO_SDK_TEST_USER and UPCLOUD_GO_SDK_TEST_PASSWORD
    ctx := context.Background()
    svc := getTestService()

    // Create test server with VNC enabled
    server := createTestServerWithVNC(ctx, t, svc)
    defer deleteTestServer(ctx, t, svc, server.UUID)

    options := &VNCOptions{
        ViewOnly: true, // Don't send input in test
    }

    // This will fail without X11 display, but validates the connection setup
    err := ConnectVNC(ctx, svc, server.UUID, options)

    // Expected to fail in CI environment without display
    if err != nil && !strings.Contains(err.Error(), "DISPLAY") {
        t.Errorf("Unexpected error: %v", err)
    }
}
```

---

## UpCloud Serial Console / Text Mode

### Current Status

Based on research, **UpCloud does not provide text-mode serial console or PTY access** through their API.

Available console options:
1. **VNC** (graphical) - Supported via API ✅
2. **HTML5 Web Console** (browser) - Available in Control Panel
3. **SSH** - Traditional text access (requires network connectivity)

### Why No Serial Console?

Serial console access is typically a KVM/hypervisor feature that requires:
- Direct PTY access to the hypervisor
- WebSocket or similar real-time protocol for streaming
- Additional security considerations for raw console access

### Workaround: SSH as Primary Text Access

For text-mode access, SSH remains the recommended approach:

```bash
upctl server ssh <uuid>
```

This could be enhanced to:
1. Automatically fetch server IP addresses
2. Configure SSH keys if not already set
3. Handle SSH connection in one command

### Future Enhancement Request

If text-mode console is needed, consider opening a feature request with UpCloud for:
- Serial console API endpoint (similar to AWS EC2 Serial Console)
- WebSocket-based text console streaming
- Read-only console log access

---

## VNC vs Serial Console Comparison

| Feature | VNC | Serial Console |
|---------|-----|----------------|
| **Type** | Graphical | Text-only |
| **Use Cases** | Full desktop, troubleshooting GUI | GRUB menu, boot logs, emergency access |
| **Network Required** | No (works even with broken network) | No |
| **UpCloud Support** | ✅ Yes (API) | ❌ No |
| **Boot Access** | Limited | Full (BIOS/GRUB) |
| **Performance** | Bandwidth-heavy | Lightweight |
| **Security** | Password + optional encryption | Usually password-only |

---

## Security Considerations

### Password Storage

**Never:**
- Store VNC passwords in configuration files
- Pass passwords via command line arguments
- Log passwords to stdout/stderr
- Keep password files after connection ends

**Always:**
- Use temporary files with 0600 permissions
- Clean up password files immediately after use
- Use stdin when possible (`-autopass`)
- Set up signal handlers for cleanup on interrupt

### Connection Security

VNC connections are **not encrypted by default**. For production use:

1. **Use SSH tunneling:**
```bash
ssh -L 5901:localhost:5901 user@server
vncviewer localhost:5901
```

2. **Enable VNC encryption** (if client supports):
```bash
vncviewer -SecurityTypes VeNCrypt,TLSVnc server:5901
```

3. **Use VPN** for accessing VNC over public internet

---

## References

### VNC Client Documentation

- [TigerVNC Official Site](https://tigervnc.org/)
- [TigerVNC vncviewer Manual](https://tigervnc.org/doc/vncviewer.html)
- [TigerVNC GitHub](https://github.com/TigerVNC/tigervnc)
- [RealVNC Viewer Parameter Reference](https://help.realvnc.com/hc/en-us/articles/360002254618-RealVNC-Viewer-Parameter-Reference)
- [Ubuntu VNC Manual Page](https://manpages.ubuntu.com/manpages/xenial/man1/xtightvncviewer.1.html)

### UpCloud Documentation

- [Connecting to Cloud Server over VNC](https://upcloud.com/docs/guides/connecting-to-cloud-server-over-vnc/)
- [UpCloud API - Servers](https://developers.upcloud.com/1.3/8-servers/)
- [Connecting to Your Server](https://upcloud.com/docs/guides/connecting-to-your-server/)

### Serial Console References (Other Providers)

- [AWS EC2 Serial Console](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-serial-console.html)
- [Google Cloud Serial Console](https://cloud.google.com/compute/docs/troubleshooting/troubleshooting-using-serial-console)
- [How to Enable virsh Console Access](https://ostechnix.com/how-to-enable-virsh-console-access-for-kvm-guests/)
- [Connecting to VMs with libvirt](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/9/html/configuring_and_managing_virtualization/assembly_connecting-to-virtual-machines_configuring-and-managing-virtualization)

---

## Summary

### Key Recommendations

1. **Implement `upctl server console <uuid>`** command
2. **Use TigerVNC as default** (cross-platform, open source)
3. **Never expose passwords in command line** - use temp files or stdin
4. **Auto-detect available VNC clients** on user's system
5. **Clean up credentials** immediately after use
6. **Document SSH** as primary text-mode access method

### Implementation Checklist

- [ ] Implement VNC client detection logic
- [ ] Create secure password file handling
- [ ] Add `upctl server console` command
- [ ] Support `--viewer`, `--fullscreen`, `--view-only` flags
- [ ] Add cleanup handlers for interrupts
- [ ] Write comprehensive tests
- [ ] Document installation of VNC clients
- [ ] Add SSH shortcut command for text access

### Estimated Effort

**VNC Console Command:** 4-6 hours
**Testing:** 2 hours
**Documentation:** 1 hour

**Total:** 1 working day

---

**Questions?** See TigerVNC documentation or examine similar implementations in other cloud CLI tools (AWS, Azure, GCP all provide console access commands).
