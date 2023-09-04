package server

import (
	"bytes"
	"github.com/egovorukhin/egolog"
	"regexp"
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

func (l *Logger) Write(data []byte) (n int, err error) {
	egolog.Infofn(l.Filename, string(l.maskPwd(data, []byte("password"), []byte("******"))))
	return len(data), nil
}

func (l *Logger) maskPwd(data, wordPwd, mask []byte) []byte {
	if bytes.Contains(data, wordPwd) {
		index := -1
		data = regexp.MustCompile(`[^"\\]+(?:\\.[^"\\]*)*`).ReplaceAllFunc(data, func(b []byte) []byte {
			if bytes.Contains(b, wordPwd) {
				index = 0
			}
			if index > -1 {
				index++
			}
			if index == 3 {
				return mask
			}
			return b
		})
	}
	return data
}
