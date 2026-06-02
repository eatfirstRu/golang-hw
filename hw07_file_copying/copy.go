package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func checkFilePath(path string) error {
	path = filepath.Dir(path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
		path = filepath.Dir(execPath) + path
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("not exists current directory + path: %w", err)
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := checkFilePath(fromPath)
	if err != nil {
		return err
	}

	fSrc, err := os.OpenFile(fromPath, os.O_RDONLY, 0o666)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer fSrc.Close()

	fi, err := fSrc.Stat()
	if err != nil {
		return fmt.Errorf("source file info: %w", err)
	}
	if offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}
	_, err = fSrc.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("source file seek to offset: %w", err)
	}
	if limit == 0 || limit > fi.Size()-offset {
		limit = fi.Size() - offset
	}

	err = checkFilePath(toPath)
	if err != nil {
		return err
	}
	fDst, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create destination file: [%v] %w", toPath, err)
	}
	defer fDst.Close()

	written, err := io.CopyN(fDst, fSrc, limit)
	if (err != nil && !errors.Is(err, io.EOF)) || written < limit {
		return fmt.Errorf("copy to destination file: %w", err)
	}

	return nil
}
