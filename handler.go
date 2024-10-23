// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
	"io"
)

type Handler struct {
	rtmp.DefaultHandler

	ConnMsg *rtmpmsg.NetConnectionConnect
	PubMsg  *rtmpmsg.NetStreamPublish

	Ctx    context.Context
	Cancel context.CancelFunc

	Pipe    *Pipe
	FlvFile *FlvFile

	HandleFunc func(h *Handler) error
}

func (h *Handler) OnConnect(timestamp uint32, cmd *rtmpmsg.NetConnectionConnect) error {
	h.Pipe = NewPipe()

	h.Ctx, h.Cancel = context.WithCancel(context.Background())
	h.ConnMsg = cmd
	return nil
}

func (h *Handler) OnPublish(_ *rtmp.StreamContext, timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error {
	h.PubMsg = cmd
	return h.HandleFunc(h)
}

func (h *Handler) OnAudio(timestamp uint32, payload io.Reader) error {
	var data flvtag.AudioData
	err := flvtag.DecodeAudioData(payload, &data)
	if err != nil {
		return err
	}
	h.Pipe.WriteAudio(timestamp, &data)

	return h.FlvFile.WriteAudio(timestamp, &data)
}

func (h *Handler) OnVideo(timestamp uint32, payload io.Reader) error {
	var data flvtag.VideoData
	err := flvtag.DecodeVideoData(payload, &data)
	if err != nil {
		return err
	}
	h.Pipe.WriteVideo(timestamp, &data)

	return h.FlvFile.WriteVideo(timestamp, &data)
}

func (h *Handler) OnClose() {
	h.Cancel()
	h.Pipe.Close()
}
