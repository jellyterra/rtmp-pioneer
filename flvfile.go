// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/yutopp/go-flv"
	flvtag "github.com/yutopp/go-flv/tag"
	"os"
	"path/filepath"
	"time"
)

type FlvFile struct {
	File *os.File
	Enc  *flv.Encoder
}

func CreateFlvFile(baseDir string) (*FlvFile, error) {
	f, err := os.OpenFile(filepath.Join(baseDir, fmt.Sprint(time.Now().Unix(), ".flv")), os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}

	enc, err := flv.NewEncoder(f, flv.FlagsAudio|flv.FlagsVideo)
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	return &FlvFile{
		File: f,
		Enc:  enc,
	}, nil
}

func (e *FlvFile) WriteAudio(timestamp uint32, data *flvtag.AudioData) error {
	return e.Enc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeAudio,
		Timestamp: timestamp,
		Data:      data,
	})
}

func (e *FlvFile) WriteVideo(timestamp uint32, data *flvtag.VideoData) error {
	return e.Enc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeVideo,
		Timestamp: timestamp,
		Data:      data,
	})
}
