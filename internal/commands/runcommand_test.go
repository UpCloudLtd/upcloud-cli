package commands

import (
	"bytes"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"testing"
)

type mockNone struct {
	*cobra.Command
	mock.Mock
}

func (m *mockNone) InitCommand() {
	panic("implement me")
}

func (m *mockNone) Cobra() *cobra.Command {
	return m.Command
}

func (m *mockNone) ExecuteWithoutArguments(exec Executor) (output.Output, error) {
	args := m.Called(exec)
	return args.Get(0).(output.Output), args.Error(1)
}

func TestRunCommand(t *testing.T) {
	var mockNoArgs = &mockNone{Command: &cobra.Command{}}
	mockNoArgs.On("ExecuteWithoutArguments", mock.Anything).Return(output.OnlyMarshaled{Value: "mock"}, nil)
	mService := &smock.Service{}
	cfg := config.New()
	cfg.Viper().Set(config.KeyOutput, config.ValueOutputJSON)
	// capture stdout
	oldStdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// run
	err := commandRunE(mockNoArgs, mService, cfg, []string{})
	assert.NoError(t, err)
	mockNoArgs.AssertNumberOfCalls(t, "ExecuteWithoutArguments", 1)

	os.Stdout = oldStdOut
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}

	// validate stdout
	// TODO: this prooobably should be made a bit nicer as well?
	assert.Equal(t, "\"mock\"\n", buf.String())
}
