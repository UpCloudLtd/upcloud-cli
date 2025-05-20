package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	internal "github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

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

func (m *mockNone) InitCommandWithConfig(*config.Config) {
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

func (m *mockSingle) InitCommandWithConfig(*config.Config) {
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

func (m *mockMulti) InitCommandWithConfig(*config.Config) {
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

func (m *mockMultiResolver) Get(_ context.Context, _ internal.AllServices) (resolver.Resolver, error) {
	return func(arg string) resolver.Resolved {
		rv := resolver.Resolved{Arg: arg}
		if len(arg) <= 5 {
			rv.AddMatch("uuid:"+arg, resolver.MatchTypeExact)
		}
		return rv
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

func (m *mockMultiResolver) InitCommandWithConfig(*config.Config) {
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
			expectedRunError: "exactly one positional argument is required",
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
			expectedRunError: "at least one positional argument is required",
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
			mService.On("GetAccount").Return(nil, nil)
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

func TestExecute_Offline(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetAccount").Return(nil, nil)
	cmd := &mockNone{Command: &cobra.Command{}}
	cmd.On("ExecuteWithoutArguments", mock.Anything).Return(output.OnlyMarshaled{Value: "mock"}, nil)

	cfg := config.New()
	cfg.Viper().Set(config.KeyOutput, config.ValueOutputJSON)

	err := commandRunE(cmd, mService, cfg, []string{})
	assert.NoError(t, err)
}

func TestExecute_Resolution(t *testing.T) {
	cmd := &mockMultiResolver{Command: &cobra.Command{}}
	mService := &smock.Service{}
	cfg := config.New()
	cfg.Viper().Set(config.KeyOutput, config.ValueOutputJSON)
	executor := NewExecutor(cfg, mService, flume.New("test"))
	outputs, err := execute(cmd, executor, []string{"a", "b", "failtoresolve", "c"}, resolveOnly, 10, func(_ Executor, arg string) (output.Output, error) {
		return output.OnlyMarshaled{Value: arg}, nil
	})
	assert.Len(t, outputs, 4)
	assert.NoError(t, err)

	// as results are run in parallel, they dont always come out in the same order.
	values := map[string]struct{}{}
	for _, o := range outputs {
		switch typedO := o.(type) {
		case output.OnlyMarshaled:
			values[typedO.Value.(string)] = struct{}{}
		case output.Error:
			assert.Empty(t, typedO.Resolved)
			assert.EqualError(t, typedO.Value, "cannot resolve argument: nothing found matching 'failtoresolve'")
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
	outputs, err := execute(cmd, executor, []string{"a", "b", "failToExecute", "c"}, resolveOnly, 10, cmd.Execute)
	assert.Len(t, outputs, 4)
	assert.NoError(t, err)

	for _, o := range outputs {
		switch typedO := o.(type) {
		case output.OnlyMarshaled:
			assert.Equal(t, output.OnlyMarshaled{Value: "mock"}, typedO)
		case output.Error:
			assert.Equal(t, "failToExecute", typedO.Original, typedO.Resolved)
			assert.EqualError(t, typedO.Value, "MOCKKFFAIL")
		}
	}
}
