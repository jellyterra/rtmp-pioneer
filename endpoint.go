// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	rtmpmsg "github.com/yutopp/go-rtmp/message"
)

type Endpoint interface {
	WriteAudio(timestamp uint32, payload []byte) error
	WriteVideo(timestamp uint32, payload []byte) error
	WriteSetFrame(timestamp uint32, data *rtmpmsg.DataMessage) error
	Close()
}
