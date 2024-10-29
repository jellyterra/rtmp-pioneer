// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
)

type Remote struct {
	Client *rtmp.ClientConn
	Stream *rtmp.Stream
}

func DialRemote(h *Handler, host, app string) (*Remote, error) {

	connMsg := h.ConnMsg.Command
	connMsg.App = app
	connMsg.TCURL = "rtmp://" + host + "/" + app

	fmt.Println(h.Time, "Connecting to", connMsg.TCURL)

	client, err := rtmp.Dial("rtmp", host+":1935", &rtmp.ConnConfig{})
	if err != nil {
		return nil, err
	}

	err = client.Connect(&rtmpmsg.NetConnectionConnect{
		Command: connMsg,
	})
	if err != nil {
		return nil, err
	}

	stream, err := client.CreateStream(&rtmpmsg.NetConnectionCreateStream{}, 128)
	if err != nil {
		return nil, err
	}

	err = stream.Publish(h.PubMsg)
	if err != nil {
		return nil, err
	}

	return &Remote{
		Client: client,
		Stream: stream,
	}, nil
}

func (e *Remote) Close() {
	_ = e.Stream.Close()
	_ = e.Client.Close()
}

func (e *Remote) WriteAudio(timestamp uint32, payload []byte) error {
	return e.Stream.Write(5, timestamp, &rtmpmsg.AudioMessage{
		Payload: bytes.NewReader(payload),
	})
}

func (e *Remote) WriteVideo(timestamp uint32, payload []byte) error {
	return e.Stream.Write(6, timestamp, &rtmpmsg.VideoMessage{
		Payload: bytes.NewReader(payload),
	})
}

func (e *Remote) WriteSetFrame(timestamp uint32, data *rtmpmsg.DataMessage) error {
	return e.Stream.Write(8, timestamp, data)
}
