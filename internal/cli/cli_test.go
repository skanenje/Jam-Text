package cli

import (
	"strings"
	"testing"
)

func TestRunValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no command",
			args:    []string{"program"},
			wantErr: true,
			errMsg:  "unknown command",
		},
		{
			name:    "ivalid command",
			args:    []string{"program", "-cmd", "invalid"},
			wantErr: true,
			errMsg:  "unknown command: invalid",
		},
		{
			name:    "index without input",
			args:    []string{"program", "-cmd", "index"},
			wantErr: true,
			errMsg:  "input and output file paths must be specified",
		},
		{
			name:    "index without output",
			args:    []string{"program", "-cmd", "index"},
			wantErr: true,
			errMsg:  "input and output file paths must be specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Run() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

