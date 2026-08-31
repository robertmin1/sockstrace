package main

import (
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestResolveProgramArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		configArgs  []string
		programArgs []string
		argsChanged bool
		want        []string
	}{
		{
			name:        "positional arguments override config arguments",
			configArgs:  []string{"from-config"},
			programArgs: []string{"from-cli"},
			argsChanged: false,
			want:        []string{"from-cli"},
		},
		{
			name:        "--args takes precedence over positional arguments",
			configArgs:  []string{"from-args-flag"},
			programArgs: []string{"from-cli"},
			argsChanged: true,
			want:        []string{"from-args-flag"},
		},
		{
			name:        "config arguments are used when no positional arguments exist",
			configArgs:  []string{"from-config"},
			programArgs: nil,
			argsChanged: false,
			want:        []string{"from-config"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := resolveProgramArgs(
				tt.configArgs,
				tt.programArgs,
				tt.argsChanged,
			)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveProgramArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCLIParsingFlow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       []string
		configArgs  []string
		wantProgram string
		wantArgs    []string
	}{
		{
			name:        "positional arguments",
			input:       []string{"wget", "google.com"},
			wantProgram: "wget",
			wantArgs:    []string{"google.com"},
		},
		{
			name:        "double dash arguments",
			input:       []string{"wget", "--", "google.com", "-I", "--debug"},
			wantProgram: "wget",
			wantArgs:    []string{"google.com", "-I", "--debug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flags.StringSlice("args", tt.configArgs, "")

			if err := flags.Parse(tt.input); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			parsedArgs := flags.Args()
			if len(parsedArgs) == 0 {
				t.Fatal("expected target program")
			}

			program := parsedArgs[0]
			programArgs := parsedArgs[1:]

			gotArgs := resolveProgramArgs(
				tt.configArgs,
				programArgs,
				flags.Changed("args"),
			)

			if program != tt.wantProgram {
				t.Errorf("program = %q, want %q", program, tt.wantProgram)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("args = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
