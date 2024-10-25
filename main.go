// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"github.com/yutopp/go-rtmp"
	"io"
	"net"
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

						ep, err := DialRemote(h)
						if err != nil {
							return err
						}

						flvFile, err := CreateFlvFile(*outDir, fmt.Sprint(h.Time))
						if err != nil {
							return err
						}

						h.Endpoints = []Endpoint{
							flvFile,
							ep,
						}

						fmt.Println(h.Time, "Streaming started.")

						return nil
					},
				},

				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},
			}
		},
	})

	fmt.Println("Listen on", *serveAddr)
	err = srv.Serve(l)
	if err != nil {
		return err
	}

	return nil
}
