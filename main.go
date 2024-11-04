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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	outDir     = flag.String("o", "./", "Stream save directory.")
	serveAddr  = flag.String("a", ":1935", "Server listen address.")
	expireDays = flag.Int("expire", 0, "Expiration days.")
)

func main() {
	flag.Parse()

	err := _main()
	if err != nil {
		fmt.Println(err)
	}
}

func _main() error {

	if *expireDays != 0 {
		go autoExpire()
	}

	l, err := net.Listen("tcp", *serveAddr)
	if err != nil {
		return err
	}

	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			return conn, &rtmp.ConnConfig{
				Handler: &Handler{
					HandleFunc: handleConn,
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

func autoExpire() {
	for {
		err := func() error {
			entries, err := os.ReadDir(*outDir)
			if err != nil {
				return err
			}

			deadline := time.Now().AddDate(0, 0, -*expireDays)

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				ut, err := strconv.ParseInt(strings.Split(entry.Name(), ".")[0], 10, 64)
				if err != nil {
					return err
				}
				t := time.UnixMicro(ut)

				if t.Before(deadline) {
					fmt.Println("Expired", t.String())
					_ = os.Remove(filepath.Join(*outDir, entry.Name()))
				}
			}

			return nil
		}()
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Hour * 1)
	}
}

func handleConn(h *Handler) error {
	sp := strings.Split(h.Url.Path, "/")

	switch sp[1] {
	case "direct":
		return handleDirect(h, sp[2], sp[3])
	case "record":
		return handleRecord(h)
	default:
		return fmt.Errorf("unexpected route: /%s", sp[1])
	}
}

func handleDirect(h *Handler, host, app string) error {

	fmt.Println(h.Time, "Direct route.")

	ep, err := DialRemote(h, host, app)
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
}

func handleRecord(h *Handler) error {

	fmt.Println(h.Time, "No route: recording only.")

	flvFile, err := CreateFlvFile(*outDir, fmt.Sprint(h.Time))
	if err != nil {
		return err
	}

	h.Endpoints = []Endpoint{
		flvFile,
	}

	fmt.Println(h.Time, "Recording started.")

	return nil
}
