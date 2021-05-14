package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/gemalto/flume"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type mockSingle struct {
	*cobra.Command
	mock.Mock
}

func (m *mockSingle) InitCommand() {
	panic("implement me")
}

func (m *mockSingle) Cobra() *cobra.Command {
	return m.Command
}

func (m *mockSingle) ExecuteSingleArgument(exec Executor, arg string) (output.Output, error) {
	args := m.Called(exec, arg)
	return args.Get(0).(output.Output), args.Error(1)
}

type mockMulti struct {
	*cobra.Command
	mock.Mock
}

func (m *mockMulti) MaximumExecutions() int {
	return 10
}

func (m *mockMulti) InitCommand() {
	panic("implement me")
}

func (m *mockMulti) Cobra() *cobra.Command {
	return m.Command
}

func (m *mockMulti) Execute(exec Executor, arg string) (output.Output, error) {
	args := m.Called(exec, arg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Output), args.Error(1)
}

type mockMultiResolver struct {
	*cobra.Command
	mock.Mock
}

func (m *mockMultiResolver) Get(_ internal.AllServices) (resolver.Resolver, error) {
	return func(arg string) (uuid string, err error) {
		if len(arg) > 5 {
			return "", fmt.Errorf("MOCKTOOLONG")
		}
		return fmt.Sprintf("uuid:%s", arg), nil
	}, nil
}

func (m *mockMultiResolver) PositionalArgumentHelp() string {
	return "MOCKARG"
}

func (m *mockMultiResolver) MaximumExecutions() int {
	return 10
}

func (m *mockMultiResolver) InitCommand() {
	panic("implement me")
}

func (m *mockMultiResolver) Cobra() *cobra.Command {
	return m.Command
}

func (m *mockMultiResolver) Execute(exec Executor, arg string) (output.Output, error) {
	args := m.Called(exec, arg)
	return args.Get(0).(output.Output), args.Error(1)
}

func TestRunCommand(t *testing.T) {
	for _, test := range []struct {
		name             string
		args             []string
		command          Command
		expectedResult   string
		expectedRunError string
	}{
		{
			name:           "no args",
			args:           []string{},
			command:        &mockNone{Command: &cobra.Command{}},
			expectedResult: "\"mock\"\n",
		},
		{
			name:           "no args with args should ignore arguments",
			args:           []string{"goo"},
			command:        &mockNone{Command: &cobra.Command{}},
			expectedResult: "\"mock\"\n",
		},
		{
			name:           "single arg",
			args:           []string{"hello"},
			command:        &mockSingle{Command: &cobra.Command{}},
			expectedResult: "\"mock\"\n",
		},
		{
			name:             "single arg with no args",
			args:             []string{},
			command:          &mockSingle{Command: &cobra.Command{}},
			expectedRunError: "exactly 1 argument is required",
		},
		{
			name:           "multi args",
			args:           []string{"foo", "bar"},
			command:        &mockMulti{Command: &cobra.Command{}},
			expectedResult: "[\n  \"mock\",\n  \"mock\"\n]\n",
		},
		{
			name:             "multiarg with no args",
			args:             []string{},
			command:          &mockMulti{Command: &cobra.Command{}},
			expectedRunError: "at least one argument is required",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			switch test.command.(type) {
			case NoArgumentCommand:
				test.command.(*mockNone).On("ExecuteWithoutArguments", mock.Anything).Return(output.OnlyMarshaled{Value: "mock"}, nil)
			case SingleArgumentCommand:
				test.command.(*mockSingle).On("ExecuteSingleArgument", mock.Anything, mock.Anything).Return(output.OnlyMarshaled{Value: "mock"}, nil)
			case MultipleArgumentCommand:
				test.command.(*mockMulti).On("Execute", mock.Anything, mock.Anything).Return(output.OnlyMarshaled{Value: "mock"}, nil)
			}
			mService := &smock.Service{}
			cfg := config.New()
			cfg.Viper().Set(config.KeyOutput, config.ValueOutputJSON)
			// capture stdout
			oldStdOut := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// run
			err := commandRunE(test.command, mService, cfg, test.args)
			if test.expectedRunError == "" {
				assert.NoError(t, err)
				switch test.command.(type) {
				case NoArgumentCommand:
					test.command.(*mockNone).AssertNumberOfCalls(t, "ExecuteWithoutArguments", 1)
				case SingleArgumentCommand:
					test.command.(*mockSingle).AssertNumberOfCalls(t, "ExecuteSingleArgument", 1)
				case MultipleArgumentCommand:
					test.command.(*mockMulti).AssertNumberOfCalls(t, "Execute", len(test.args))
				}
			} else {
				assert.EqualError(t, err, test.expectedRunError)
			}

			os.Stdout = oldStdOut
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}

			if test.expectedResult != "" {
				var buf bytes.Buffer
				if _, err := io.Copy(&buf, r); err != nil {
					t.Fatal(err)
				}
				// validate stdout
				assert.Equal(t, test.expectedResult, buf.String())
			}
		})
	}
}

func TestExecute_Resolution(t *testing.T) {
	cmd := &mockMultiResolver{Command: &cobra.Command{}}
	mService := &smock.Service{}
	cfg := config.New()
	cfg.Viper().Set(config.KeyOutput, config.ValueOutputJSON)
	executor := NewExecutor(cfg, mService, flume.New("test"))
	results, err := execute(cmd, executor, []string{"a", "b", "failtoresolve", "c"}, 10, func(exec Executor, arg string) (output.Output, error) {
		return output.OnlyMarshaled{Value: arg}, nil
	})
	assert.Len(t, results, 4)
	assert.NoError(t, err)

	// as results are run in parallel, they dont always come out in the same order.
	values := map[string]struct{}{}
	for _, r := range results {
		if r.Error == nil {
			values[r.Result.(output.OnlyMarshaled).Value.(string)] = struct{}{}
		} else {
			assert.Nil(t, r.Result)
			assert.EqualError(t, r.Error, "cannot resolve argument: MOCKTOOLONG")
		}
	}
	assert.Equal(t, values, map[string]struct{}{
		"uuid:a": {},
		"uuid:b": {},
		"uuid:c": {},
	})
}

func TestExecute_Error(t *testing.T) {
	cmd := &mockMulti{Command: &cobra.Command{}}
	cmd.On("Execute", mock.Anything, mock.MatchedBy(func(arg string) bool {
		return len(arg) < 5
	})).Return(output.OnlyMarshaled{Value: "mock"}, nil)
	cmd.On("Execute", mock.Anything, mock.MatchedBy(func(arg string) bool {
		return len(arg) >= 5
	})).Return(nil, fmt.Errorf("MOCKKFFAIL"))
	mService := &smock.Service{}
	cfg := config.New()
	cfg.Viper().Set(config.KeyOutput, config.ValueOutputJSON)
	executor := NewExecutor(cfg, mService, flume.New("test"))
	results, err := execute(cmd, executor, []string{"a", "b", "failtorexecute", "c"}, 10, cmd.Execute)
	assert.Len(t, results, 4)
	assert.NoError(t, err)

	for _, r := range results {
		if r.Error == nil {
			assert.Equal(t, output.OnlyMarshaled{Value: "mock"}, r.Result)
		} else {
			assert.Nil(t, r.Result)
			assert.EqualError(t, r.Error, "MOCKKFFAIL")
		}
	}
}
