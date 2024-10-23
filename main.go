// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/yutopp/go-rtmp"
)

var (
	outDir    = flag.String("o", "./", "Stream save directory.")
	serveAddr = flag.String("a", ":1935", "Server listen address.")
)

func main() {
	flag.Parse()

	err := _main()
	if err != nil {
		fmt.Println(err)
	}
}

func _main() error {

	l, err := net.Listen("tcp", *serveAddr)
	if err != nil {
		return err
	}

	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			return conn, &rtmp.ConnConfig{
				Handler: &Handler{
					HandleFunc: func(h *Handler) error {
						h.FlvFile, err = CreateFlvFile(*outDir)
						if err != nil {
							return err
						}

						return nil
					},
				},
				Logger: log.StandardLogger(),

				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},
			}
		},
	})

	err = srv.Serve(l)
	if err != nil {
		return err
	}

	return nil
}
