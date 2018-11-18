package utils

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// hard to test - would need more abstraction
func Gzip(source string) error {
	i, err := os.Open(source)
	defer i.Close()
	if err != nil {
		return err
	}

	output := withoutDotEnding(source) + ".zip"
	o, err := os.Create(output)
	defer o.Close()
	if err != nil {
		return err
	}

	gw := gzip.NewWriter(o)
	gw.Name = filepath.Base(source)
	defer gw.Close()
	io.Copy(gw, i)

	return nil
}

func withoutDotEnding(src string) string {
	runes := []rune(src)

	for i := len(runes) - 1; i > 0; i-- {
		if runes[i] == '.' {
			return string(runes[:i])
		}
	}
	return src
}
