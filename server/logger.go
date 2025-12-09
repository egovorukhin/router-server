package server

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/egovorukhin/egolog"
)

type Logger struct {
	Name     string `yaml:"name"`
	Format   string `yaml:"format"`
	Filename string `yaml:"filename"`
	Time     Time   `yaml:"time"`
}

type Time struct {
	Format   string `yaml:"time_format"`
	Zone     string `yaml:"time_zone"`
	Interval int    `yaml:"time_interval"`
}

type Loggers []Logger

func (l Loggers) Get(name string) *Logger {
	for _, log := range l {
		if log.Name == name {
			return &log
		}
	}
	return nil
}

var (
	pwd  = []byte("password")
	mask = []byte("******")
)

func (l *Logger) Write(data []byte) (n int, err error) {
	return len(l.convert(data)), nil
}

func (l *Logger) convert(data []byte) []byte {
	switch {
	case strings.Contains(l.Format, "${body}") && bytes.Contains(data, []byte("/api/v1/auth")):
		if matches := regexp.MustCompile(`\{[[:graph:]\s]*?}`).FindSubmatch(data); len(matches) > 0 {
			is := false
			for _, field := range regexp.MustCompile(`[^"\\]+(?:\\.[^"\\]*)*`).
				FindAll(bytes.ToLower(l.trim(matches[0])), -1) {
				if bytes.Equal(field, []byte("{")) || bytes.Equal(field, []byte(":")) ||
					bytes.Equal(field, []byte(",")) || bytes.Equal(field, []byte("}")) {
					continue
				}
				if is {
					data = bytes.ReplaceAll(data, field, mask)
					continue
				}
				if bytes.Equal(field, pwd) {
					is = true
					continue
				}
				is = false
			}
		}
	}

	l.Info(string(data))
	return data
}

func (l *Logger) trim(s []byte) []byte {
	return bytes.Map(func(r rune) rune {
		switch r {
		case 0x0009, 0x000A, 0x0020, 0x000B, 0x000C, 0x000D, 0x0085, 0x2028, 0x2029:
			return -1
		default:
			return r
		}
	}, s)
}

func (l *Logger) Info(format string, v ...interface{}) {
	egolog.Infofn(l.Filename, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	egolog.Errorfn(l.Filename, format, v...)
}
