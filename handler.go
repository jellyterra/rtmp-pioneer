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
	"net/url"
	"time"
)

type Handler struct {
	rtmp.DefaultHandler

	ConnMsg *rtmpmsg.NetConnectionConnect
	PubMsg  *rtmpmsg.NetStreamPublish
	Url     *url.URL

	Time int64

	Endpoints []Endpoint

	HandleFunc func(h *Handler) error

	HooksOnClose []func()
}

func (h *Handler) Logln(a ...interface{}) {
	fmt.Println(append([]any{h.Time}, a...)...)
}

func (h *Handler) Close() {
	for _, ep := range h.Endpoints {
		ep.Close()
	}
}

func (h *Handler) EndpointError(i int, err error) error {
	h.Logln("Endpoint", i, "error:", err)
	return err
}

func (h *Handler) OnConnect(timestamp uint32, cmd *rtmpmsg.NetConnectionConnect) (err error) {
	h.ConnMsg = cmd
	h.Url, err = url.Parse(cmd.Command.TCURL)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) OnPublish(_ *rtmp.StreamContext, _ uint32, cmd *rtmpmsg.NetStreamPublish) error {
	h.PubMsg = cmd
	h.Time = time.Now().UnixMicro()

	err := h.HandleFunc(h)
	if err != nil {
		h.Logln(err)
		return err
	}

	return err
}

func (h *Handler) OnAudio(timestamp uint32, payload io.Reader) error {
	p, err := io.ReadAll(payload)
	if err != nil {
		return err
	}

	for i, ep := range h.Endpoints {
		err := ep.WriteAudio(timestamp, p)
		if err != nil {
			return h.EndpointError(i, err)
		}
	}

	return nil
}

func (h *Handler) OnVideo(timestamp uint32, payload io.Reader) error {
	p, err := io.ReadAll(payload)
	if err != nil {
		return err
	}

	for i, ep := range h.Endpoints {
		err := ep.WriteVideo(timestamp, p)
		if err != nil {
			return h.EndpointError(i, err)
		}
	}

	return nil
}

func (h *Handler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	for i, ep := range h.Endpoints {
		err := ep.WriteSetFrame(timestamp, &rtmpmsg.DataMessage{
			Name:     "@setDataFrame",
			Encoding: rtmpmsg.EncodingTypeAMF0,
			Body:     bytes.NewReader(data.Payload),
		})
		if err != nil {
			return h.EndpointError(i, err)
		}
	}

	return nil
}

func (h *Handler) OnClose() {
	h.Close()
	h.Logln("Closed.")
	for _, hook := range h.HooksOnClose {
		hook()
	}
}
