// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	log "github.com/sirupsen/logrus"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
)

type Endpoint struct {
	Client *rtmp.ClientConn
	Stream *rtmp.Stream
}

func NewEndpoint(addr, url, pub, app string) (*Endpoint, error) {
	client, err := rtmp.Dial("rtmp", addr, &rtmp.ConnConfig{
		Logger: log.StandardLogger(),
	})
	if err != nil {
		return nil, err
	}

	err = client.Connect(&rtmpmsg.NetConnectionConnect{
		Command: rtmpmsg.NetConnectionConnectCommand{
			App:   app,
			TCURL: url,
			Type:  "nonprivate",
		},
	})
	if err != nil {
		return nil, err
	}

	stream, err := client.CreateStream(nil, 128)
	if err != nil {
		return nil, err
	}

	err = stream.Publish(&rtmpmsg.NetStreamPublish{
		PublishingName: pub,
		PublishingType: "live",
	})
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		Client: client,
		Stream: stream,
	}, nil
}

func (e *Endpoint) Close() {
	_ = e.Stream.Close()
	_ = e.Client.Close()
}

func (e *Endpoint) WriteAudio(timestamp uint32, data *flvtag.AudioData) error {
	return e.Stream.Write(6, timestamp, &rtmpmsg.AudioMessage{
		Payload: data.Data,
	})
}

func (e *Endpoint) WriteVideo(timestamp uint32, data *flvtag.VideoData) error {
	return e.Stream.Write(5, timestamp, &rtmpmsg.VideoMessage{
		Payload: data.Data,
	})
}
