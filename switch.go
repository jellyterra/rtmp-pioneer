// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import flvtag "github.com/yutopp/go-flv/tag"

type Audio struct {
	Timestamp uint32
	Data      *flvtag.AudioData
}

type Video struct {
	Timestamp uint32
	Data      *flvtag.VideoData
}

type Pipe struct {
	AudioChan chan Audio
	VideoChan chan Video
}

func NewPipe() *Pipe {
	return &Pipe{
		AudioChan: make(chan Audio),
		VideoChan: make(chan Video),
	}
}

func (p *Pipe) WriteAudio(timestamp uint32, data *flvtag.AudioData) {
	select {
	case p.AudioChan <- Audio{
		Timestamp: timestamp,
		Data:      data,
	}:
	default:
		// Drop.
	}
}

func (p *Pipe) WriteVideo(timestamp uint32, data *flvtag.VideoData) {
	select {
	case p.VideoChan <- Video{
		Timestamp: timestamp,
		Data:      data,
	}:
	default:
		// Drop
	}
}

func (p *Pipe) Close() {
	close(p.AudioChan)
	close(p.VideoChan)
}
