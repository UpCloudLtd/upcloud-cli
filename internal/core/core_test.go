package core

import (
	"bytes"
	"context"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestInputValidation(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "prefix-")
	if err != nil {
		t.Fatal("Cannot create temporary file", err)
	}

	defer os.Remove(tmpFile.Name())

	for _, test := range []struct {
		name         string
		args         []string
		setupFixture func()
		unsetFixture func()
		error        string
		errorWindows string
	}{
		{
			name: "validate output flag",
			args: []string{
				"--output", "toml",
				"version",
			},
			error: "OutputFormat: toml does not validate as in(human|json|yaml)",
		},
		{
			name: "validate config flag",
			args: []string{
				"--config", "/invalid/file/path",
				"version",
			},
			error:        "cannot load configuration: unable to parse config from file '/invalid/file/path': open /invalid/file/path: no such file or directory",
			errorWindows: "cannot load configuration: unable to parse config from file '/invalid/file/path': open /invalid/file/path: The system cannot find the path specified.",
		},
		{
			name: "validate output config via env var",
			args: []string{
				"version",
			},
			setupFixture: func() {
				os.Setenv("UPCLOUD_OUTPUT", "toml")
			},
			unsetFixture: func() {
				os.Unsetenv("UPCLOUD_OUTPUT")
			},
			error: "output format 'toml' not accepted",
		},
		/*
			TODO: re-enable when we have a clear way of testing this in a configured environment
			{
				name: "validate no creds",
				args: []string{
					"version",
				},
				error: fmt.Sprintf("cannot create service: user credentials not found, these must be set in config file (default location %s) or via environment variables", filepath.Join(xdg.ConfigHome, "upctl.yaml")),
			},
		*/
		{
			name: "validate set credentials via env vars",
			args: []string{
				"version",
			},
			setupFixture: func() {
				os.Setenv("UPCLOUD_USERNAME", "foo_user")
				os.Setenv("UPCLOUD_PASSWORD", "foo_passwd")
			},
			unsetFixture: func() {
				os.Unsetenv("UPCLOUD_USERNAME")
				os.Unsetenv("UPCLOUD_PASSWORD")
			},
		},
		{
			name: "validate set credentials via config file",
			args: []string{
				"--config", tmpFile.Name(),
				"version",
			},
			setupFixture: func() {
				// Example writing to the file
				text := []byte("username: foo_user\npassword: foo_passwd")
				if _, err = tmpFile.Write(text); err != nil {
					t.Fatal("Failed to write to temporary file", err)
				}
			},
			unsetFixture: func() {
				// Close the file
				if err := tmpFile.Close(); err != nil {
					t.Fatal(err)
				}
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cmd := BuildCLI()
			stdout := bytes.NewBufferString("")

			if test.setupFixture != nil {
				test.setupFixture()
			}

			cmd.SetOutput(stdout) // prevent noisy prints
			cmd.SetArgs(test.args)

			err := cmd.Execute()

			if test.error != "" {
				if err == nil {
					t.Fatalf("expected error '%v', got no error", test.error)
				}

				if runtime.GOOS == "windows" && test.errorWindows != "" {
					assert.Equal(t, test.errorWindows, err.Error())
				} else {
					assert.Equal(t, test.error, err.Error())
				}
			}

			if test.unsetFixture != nil {
				test.unsetFixture()
			}
		})
	}
}

func TestExecute(t *testing.T) {
	realArgs := os.Args
	defer func() { os.Args = realArgs }()

	for _, test := range []struct {
		name     string
		args     []string
		expected int
	}{
		{
			name:     "Successful command (upctl version)",
			args:     []string{"upctl", "version"},
			expected: 0,
		},
		{
			name:     "Failing command (upctl server create)",
			args:     []string{"upctl", "server", "create"},
			expected: 100,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.args
			assert.Equal(t, test.expected, Execute())
		})
	}
}

func TestConfigContextAndCleanup(t *testing.T) {
	t.Run("BuildCLIWithContext should cleanup on context cancellation", func(t *testing.T) {
		// Create parent context with cancel
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Build CLI with context
		rootCmd := BuildCLIWithContext(ctx)
		rootCmd.SetArgs([]string{"version"})

		// Start command execution in background
		done := make(chan struct{})
		go func() {
			defer close(done)
			err := rootCmd.Execute()
			assert.NoError(t, err)
		}()

		// Cancel context and verify cleanup
		cancel()

		// Wait for command to finish
		select {
		case <-done:
			// Success: command completed after cancellation
		case <-time.After(500 * time.Millisecond):
			t.Error("Command didn't complete after context cancellation")
		}
	})

	t.Run("Config should implement io.Closer", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		conf := config.NewWithContext(ctx)

		// Register cleanup
		cleanup := registerForCleanup(conf)

		// Run cleanup
		cleanup()

		// Verify context was cancelled (config.Close() should have been called)
		select {
		case <-conf.Context().Done():
			// Success: context was cancelled by Close()
		default:
			t.Error("Context wasn't cancelled after cleanup")
		}
	})
}
