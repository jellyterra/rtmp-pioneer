// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"github.com/yutopp/go-flv"
	flvtag "github.com/yutopp/go-flv/tag"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
	"os"
	"path/filepath"
)

type FlvFile struct {
	File *os.File
	Enc  *flv.Encoder
}

func CreateFlvFile(baseDir, filename string) (*FlvFile, error) {
	f, err := os.OpenFile(filepath.Join(baseDir, filename+".flv"), os.O_WRONLY|os.O_CREATE, 0660)
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

func (e *FlvFile) Close() {
	_ = e.File.Close()
}

func (e *FlvFile) WriteAudio(timestamp uint32, payload []byte) error {
	var data flvtag.AudioData
	err := flvtag.DecodeAudioData(bytes.NewReader(payload), &data)
	if err != nil {
		return err
	}

	return e.Enc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeAudio,
		Timestamp: timestamp,
		Data:      &data,
	})
}

func (e *FlvFile) WriteVideo(timestamp uint32, payload []byte) error {
	var data flvtag.VideoData
	err := flvtag.DecodeVideoData(bytes.NewReader(payload), &data)
	if err != nil {
		return err
	}

	return e.Enc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeVideo,
		Timestamp: timestamp,
		Data:      &data,
	})
}

func (e *FlvFile) WriteSetFrame(timestamp uint32, data *rtmpmsg.DataMessage) error {
	return nil
}
