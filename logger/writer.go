package logger

import (
	"encoding/json"
	"os"
	"regexp"
)

const (
	JSON = "json"
	TEXT = "text"
)

type Writer interface {
	Write(*Log) error
}

type FileWriterOpts struct {
	FileName string
	Format   string
}

type ToFileWriter struct {
	FileName string
	Format   string
}

func WithToFileWriter(fileName string, format string) *ToFileWriter {
	return &ToFileWriter{
		FileName: fileName,
		Format:   format,
	}
}

func (w *ToFileWriter) Write(l *Log) error {
	var b []byte
	var err error

	switch w.Format {
	case JSON:
		b, err = json.Marshal(l)
		if err != nil {
			return err
		}
		b = append(b, '\n')

	case TEXT:
		b = []byte(l.Str)
		ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
		b = ansi.ReplaceAll(b, []byte{})

	default:
		b, err = json.Marshal(l)
		if err != nil {
			return err
		}
		b = append(b, '\n')
	}

	f, err := os.OpenFile(w.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}

type ToStdoutWriter struct{}

func WithToStdoutWriter() *ToStdoutWriter {
	return &ToStdoutWriter{}
}

func (w *ToStdoutWriter) Write(l *Log) error {
	os.Stdout.Write([]byte(l.Str))
	return nil
}
