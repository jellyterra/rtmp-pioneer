// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
	"io"
	"time"
)

type Handler struct {
	rtmp.DefaultHandler

	ConnMsg *rtmpmsg.NetConnectionConnect
	PubMsg  *rtmpmsg.NetStreamPublish

	Time int64

	Endpoints []Endpoint

	HandleFunc func(h *Handler) error
}

func (h *Handler) OnConnect(timestamp uint32, cmd *rtmpmsg.NetConnectionConnect) error {
	h.ConnMsg = cmd
	return nil
}

func (h *Handler) OnPublish(_ *rtmp.StreamContext, timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error {
	h.PubMsg = cmd
	h.Time = time.Now().UnixMicro()

	err := h.HandleFunc(h)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (h *Handler) OnAudio(timestamp uint32, payload io.Reader) error {
	p, err := io.ReadAll(payload)
	if err != nil {
		return err
	}

	for _, ep := range h.Endpoints {
		err := ep.WriteAudio(timestamp, p)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (h *Handler) OnVideo(timestamp uint32, payload io.Reader) error {
	p, err := io.ReadAll(payload)
	if err != nil {
		return err
	}

	for _, ep := range h.Endpoints {
		err := ep.WriteVideo(timestamp, p)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (h *Handler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	for _, ep := range h.Endpoints {
		err := ep.WriteSetFrame(timestamp, &rtmpmsg.DataMessage{
			Name:     "@setDataFrame",
			Encoding: rtmpmsg.EncodingTypeAMF0,
			Body:     bytes.NewReader(data.Payload),
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (h *Handler) OnClose() {
	for _, ep := range h.Endpoints {
		ep.Close()
	}

	fmt.Println(h.Time, "Closed.")
}
