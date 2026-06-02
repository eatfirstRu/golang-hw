package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type testFile struct {
	from, to      string
	limit, offset int64
}

func TestCopy(t *testing.T) {
	// Place your code here.
	err := os.MkdirAll("tmp", 0o755)
	require.NoError(t, err)

	testFiles := []testFile{
		{from: "testdata/input.txt", to: "tmp/out_offset0_limit0.txt", offset: 0, limit: 0},
		{from: "testdata/input.txt", to: "tmp/out_offset0_limit10.txt", offset: 0, limit: 10},
		{from: "testdata/input.txt", to: "tmp/out_offset0_limit1000.txt", offset: 0, limit: 1000},
		{from: "testdata/input.txt", to: "tmp/out_offset0_limit10000.txt", offset: 0, limit: 10000},
		{from: "testdata/input.txt", to: "tmp/out_offset100_limit1000.txt", offset: 100, limit: 1000},
		{from: "testdata/input.txt", to: "tmp/out_offset6000_limit1000.txt", offset: 6000, limit: 1000},
	}

	t.Run("no errors during copy", func(t *testing.T) {
		for i, tf := range testFiles {
			err := Copy(tf.from, tf.to, tf.offset, tf.limit)
			tstNum := fmt.Sprintf("test #%d", i)
			require.NoError(t, err, tstNum)
		}
	})

	t.Run("equal copy", func(t *testing.T) {
		for i, tf := range testFiles {
			var bSrc bytes.Buffer
			var bDst bytes.Buffer

			fName := fmt.Sprintf("%s/out_offset%d_limit%d.txt", filepath.Dir(tf.from), tf.offset, tf.limit)
			fSrc, err := os.OpenFile(fName, os.O_RDONLY, 0o666)
			if err != nil {
				require.NoError(t, err, "open src file")
			}
			defer fSrc.Close()
			_, _ = io.Copy(&bSrc, fSrc)

			/*fSrc, err := os.OpenFile(tf.from, os.O_RDONLY, 0o666)
			if err != nil {
				require.NoError(t, err, "open src file")
			}
			defer fSrc.Close()
			fSrc.Seek(tf.offset, 0)
			fi, _ := fSrc.Stat()
			if tf.limit == 0 || tf.limit > fi.Size()-tf.offset {
				tf.limit = fi.Size() - tf.offset
			}
			_, _ = io.CopyN(&bSrc, fSrc, tf.limit)*/

			fDst, err := os.OpenFile(tf.to, os.O_RDONLY, 0o666)
			if err != nil {
				require.NoError(t, err, "open dst file")
			}
			defer fDst.Close()
			_, _ = io.Copy(&bDst, fDst)

			tstNum := fmt.Sprintf("test #%d", i)

			require.Equal(t, bSrc.String(), bDst.String(), tstNum)
		}
	})

	t.Run("error ErrOffsetExceedsFileSize", func(t *testing.T) {
		for i, tf := range testFiles {
			fi, _ := os.Stat(tf.from)
			err := Copy(tf.from, tf.to, fi.Size()*2, tf.limit)
			tstNum := fmt.Sprintf("test #%d", i)
			require.ErrorIs(t, err, ErrOffsetExceedsFileSize, tstNum)
		}
	})
}
