// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

type StreamServerConfig struct {
	Host string `json:"host"`
	Path string `json:"path"`
	Key  string `json:"key"`
}

type Profile struct {
	Remotes   []StreamServerConfig `json:"remotes"`
	Recording bool                 `json:"recording"`

	// Events:
	// beforeConnect
	// afterStart
	// afterClose
	Webhooks map[string][]*WebhookConfig `json:"webhooks"`
}
