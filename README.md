# rtmp-pioneer
Record and forward RTMP streams to multiple endpoints.

## Usage

```
Usage of rtmp-pioneer:
  -a string
        Server listen address. (default ":1935")
  -o string
        Stream save directory. (default "./")
```

Server: `rtmp://<Pioneer Addr>/<Server Addr>/<App>`

Stream key: **AS IS**

### Example

`rtmp://k1-i.jellyterra.com/live-push.bilivideo.com/live-bvc/`

Live server: `live-push.bilivideo.com/live-bvc/`

Pioneer: `k1-i.jellyterra.com`

```
$ ./rtmp-pioneer -a :1935 -o ~/Videos
INFO[0000] Listen on :1935
INFO[0005] Connecting to live-push.bilivideo.com/live-bvc
INFO[0005] 1729790878 Streaming started.
INFO[0010] 1729790878 Closed.
```
