// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type StreamServerConfig struct {
	URL string `yaml:"url"`
	Key string `yaml:"key"`
}

type Profile struct {
	Path string `yaml:"path"`

	Downstream []StreamServerConfig `yaml:"downstream"`
}

type Config struct {
	Profiles []Profile `yaml:"profiles"`
}

func ParseConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	return &config, yaml.Unmarshal(b, &config)
}
