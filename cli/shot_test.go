package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseShebang(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr error
	}{
		{
			name:  "classic python",
			input: "#!/usr/bin/python3",
			want:  []string{"/usr/bin/python3"},
		},
		{
			name:  "python with flag",
			input: "#!/usr/bin/python -u",
			want:  []string{"/usr/bin/python", "-u"},
		},
		{
			name:  "env + python3",
			input: "#!/usr/bin/env python3",
			want:  []string{"/usr/bin/env", "python3"},
		},
		{
			name:  "env -S modern style",
			input: "#!/usr/bin/env -S bash -euo pipefail",
			want:  []string{"/usr/bin/env", "-S", "bash", "-euo", "pipefail"},
		},
		{
			name:  "env -S with extra spaces",
			input: "#!  /usr/bin/env  -S   zsh -e -x  ",
			want:  []string{"/usr/bin/env", "-S", "zsh", "-e", "-x"},
		},
		{
			name:  "leading spaces before #!",
			input: "   \t  #!/bin/sh -e",
			want:  []string{"/bin/sh", "-e"},
		},
		{
			name:  "minimal shebang",
			input: "#!python3",
			want:  []string{"python3"},
		},
		{
			name:    "no #! prefix",
			input:   "/bin/sh",
			want:    nil,
			wantErr: ErrInvalidShebang,
		},
		{
			name:    "plain comment",
			input:   "# hello world",
			want:    nil,
			wantErr: ErrInvalidShebang,
		},
		{
			name:    "empty line",
			input:   "",
			want:    nil,
			wantErr: ErrInvalidShebang,
		},
		{
			name:    "only whitespace",
			input:   "   \t\n",
			want:    nil,
			wantErr: ErrInvalidShebang,
		},
		{
			name:    "just #!",
			input:   "#!",
			want:    nil,
			wantErr: ErrInvalidShebang,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseShebang(tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				if got != nil {
					t.Errorf("expected nil result on error, got %v", got)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseShebang(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
