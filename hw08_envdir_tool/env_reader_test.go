package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var TestEnv Environment = Environment{
	"BAR":   EnvValue{NeedRemove: false, Value: "bar"},
	"EMPTY": EnvValue{NeedRemove: false, Value: ""},
	"FOO":   EnvValue{NeedRemove: false, Value: "   foo\nwith new line"},
	"HELLO": EnvValue{NeedRemove: false, Value: "\"hello\""},
	"UNSET": EnvValue{NeedRemove: true, Value: ""},
}

func TestEnvironment_String(t *testing.T) {
	tests := []struct {
		name string
		e    Environment
		want string
	}{
		{
			name: "test String()",
			e:    TestEnv,
			want: "\nkey: BAR\tvalue: false,[bar]\nkey: EMPTY\tvalue: false,[]\nkey: FOO\tvalue: false,[   foo\nwith new line]\nkey: HELLO\tvalue: false,[\"hello\"]\nkey: UNSET\tvalue: true,[]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Environment.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadDir(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		want    Environment
		wantErr bool
	}{
		{
			name:    "read from bad path",
			dir:     "",
			want:    TestEnv,
			wantErr: true,
		},
		{
			name:    "read from testdata",
			dir:     "testdata/env",
			want:    TestEnv,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ReadDir(tt.dir)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ReadDir() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReadDir() succeeded unexpectedly")
			}
			if true {
				require.EqualValues(t, tt.want.String(), got.String())
			}
		})
	}
}
