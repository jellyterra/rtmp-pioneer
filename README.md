# rtmp-pioneer
Record and forward RTMP streams to multiple endpoints.

## Install

### Releases

Download the executable in [releases](https://github.com/jellyterra/rtmp-pioneer/releases).

### Build from source

```shell
go install github.com/jellyterra/rtmp-pioneer@latest
```

## Usage

```
Usage of ./rtmp-pioneer:
  -a string
        Server listen address. (default ":1935")
  -expire int
        Expiration days.
  -o string
        Stream save directory. (default "./")
```

### Expiration

Outdated files will be automatically removed.

# Route

## Direct

Server: `rtmp://<Pioneer Addr>/direct/<Server Addr>/<App>`

Stream key: **AS IS**

### Example

`rtmp://k1-i.jellyterra.com/direct/live-push.bilivideo.com/live-bvc`

Live server: `live-push.bilivideo.com/live-bvc`

Pioneer: `k1-i.jellyterra.com`

```
$ ./rtmp-pioneer -o ~/Videos
Listen on :1935
1730107018200000 Direct route.
1730107018200000 Connecting to live-push.bilivideo.com/live-bvc
1730107018200000 Streaming started.
1730107018200000 Closed.
```
