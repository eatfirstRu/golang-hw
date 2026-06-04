package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Environment map[string]EnvValue

func (e Environment) String() string {
	str := strings.Builder{}
	str.WriteString("\n")

	sl := make([]string, 0, len(e))
	for k := range e {
		sl = append(sl, k)
	}
	slices.Sort(sl)

	for _, v := range sl {
		str.WriteString(fmt.Sprintf("key: %s\tvalue: %v,[%s]\n", v, e[v].NeedRemove, e[v].Value))
	}
	return str.String()
}

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	env := make(Environment)

	// dir = filepath.Dir(dir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("get current directory: %w", err)
		}
		dir = filepath.Join(wd, dir)
	}
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("not exists current directory + path: %w", err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read files in dir %s: %w", dir, err)
	}

	for _, f := range files {
		fSrc, err := os.OpenFile(filepath.Join(dir, f.Name()), os.O_RDONLY, 0o666)
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", f.Name(), err)
		}
		defer fSrc.Close()
		reader := bufio.NewReader(fSrc)
		line, _, err := reader.ReadLine()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("read first line in file %s: %w", f.Name(), err)
		}
		str := strings.TrimRight(string(line), " \\t")
		str = strings.ReplaceAll(str, "\x00", "\n")
		env[f.Name()] = EnvValue{
			Value:      str,
			NeedRemove: errors.Is(err, io.EOF),
		}
		// fmt.Printf("f.Name: %v, str:[%v], len(str): %d \n", f.Name(), str, len(str)).
	}
	return env, nil
}
