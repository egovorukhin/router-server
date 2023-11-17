package server

import (
	"regexp"
	"strings"
)

type Router struct {
	Pattern            string   `yaml:"pattern"`
	Index              string   `yaml:"index"`
	IpClientHeaderName string   `yaml:"ipClientHeaderName"`
	Location           Location `yaml:"location"`
}

type Location map[string]struct {
	Url    string `yaml:"url"`
	Remove bool   `yaml:"remove"`
}

func (r Router) GetLocation(url string) (string, error) {
	rgx, err := regexp.Compile(r.Pattern)
	if err != nil {
		return "", err
	}
	match := rgx.FindAllString(url, -1)
	if len(match) > 0 {
		if location, ok := r.Location[match[0]]; ok {
			if location.Remove {
				url = strings.Replace(url, match[0], location.Url, 1)
			}
			return location.Url + url, nil
		}
	}
	return "", err
}
