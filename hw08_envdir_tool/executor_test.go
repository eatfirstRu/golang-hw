package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name string
		cmd  []string
		env  Environment
		want int
	}{
		{
			name: "run command",
			cmd:  []string{"command", "arg1", "arg2"},
			env:  TestEnv,
			want: 101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RunCmd(tt.cmd, tt.env)

			require.EqualValues(t, tt.want, got)

			for k, v := range tt.env {
				val := ""
				if !v.NeedRemove {
					val = v.Value
				}
				gev := os.Getenv(k)
				ln := fmt.Sprintf("env variables expected: %s=[%s], got: %s=[%s]", k, gev, k, val)
				require.Equal(t, val, gev, ln)
			}
		})
	}
}
