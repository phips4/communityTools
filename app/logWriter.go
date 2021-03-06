package app

import (
	"fmt"
	"github.com/phips4/communityTools/app/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogWriter struct {
	file *os.File
}

func NewLogWriter() *LogWriter {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err = os.Mkdir("logs", os.ModePerm)

		if err != nil {
			panic(err)
		}
	}

	fpath := filepath.Join("logs", time.Now().Format("2006-01-02")+".log")

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		_, err := os.Create(fpath)
		if err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile(fpath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	return &LogWriter{
		file: file,
	}
}

func (w *LogWriter) Write(bytes []byte) (n int, err error) {
	logLine := fmt.Sprintf("[LOG][%s]: %s", time.Now().Format("2006-01-02 15:04.05"), string(bytes))

	var printBytes int

	if printBytes, printErr := fmt.Print(logLine); printErr != nil {
		w.file.WriteString(logLine) //try write to file, ignore errors
		return printBytes, printErr
	}

	if fileBytes, fileErr := w.file.WriteString(logLine); fileErr != nil {
		return fileBytes, fileErr
	}

	// printing to the std-out has more importance than writing it to the file
	return printBytes, nil
}

func (w *LogWriter) Close() error {
	return w.file.Close()
}

func (w *LogWriter) Compress() error {
	files, err := ioutil.ReadDir("logs")
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02") + ".log"

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".log") && f.Name() != now {
			err = utils.Gzip(filepath.Join("logs", f.Name()))
			if err != nil {
				return err
			}

			err = os.Remove(filepath.Join("logs", f.Name()))
			if err != nil {
				return err
			}
		}
	}

	return err
}
