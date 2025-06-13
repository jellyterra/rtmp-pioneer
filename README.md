# RTMP Pioneer

Record and forward RTMP streams to multiple live streaming servers.

## Install

### [Releases](https://github.com/jellyterra/rtmp-pioneer/releases) download

### Build from source

```shell
go install github.com/jellyterra/rtmp-pioneer@latest
```

## Usage

```
Usage of ./rtmp-pioneer:
  -a string
        Server listening address. (default ":1935")
  -expire int
        Expiration days.
  -o string
        Streaming save directory. (default "./rec")
  -p string
        Profile directory. (default "./profile")
```

### Expiration

Outdated files will be automatically removed.

Set option `--expire` to non-zero value to enable it.

# Route

## No route (recording-only)

Server: `rtmp://<Pioneer Addr>/record`

Streaming key: **ignored**

### Example

Server: `rtmp://rtmp-pioneer.jellyterra.com/record`

```
$ rtmp-pioneer -o ~/Videos
Listen on :1935
1749123456123456 Inbound - recording only.
1749123456123456 Streaming started.
1749123456123456 Closed.
```

## Direct

Server: `rtmp://<Pioneer Addr>/direct/<Server Addr>/<Path>`

Streaming key: **AS IS to the remote server**

### Example

Server: `rtmp://rtmp-pioneer.jellyterra.com/direct/live-push.bilivideo.com/live-bvc/`

Streaming key: `?streamname=&key=&schedule=rtmp&pflag=1`

```
$ rtmp-pioneer -o ~/Videos
Listen on :1935
1749123456123456 Inbound - direct forwarding.
1749123456123456 Endpoint connecting: live-push.bilivideo.com
1749123456123456 Streaming started.
1749123456123456 Closed.
```

## Profile

Server: `rtmp://<Pioneer Addr>/profile/<Profile Name>`

Streaming key: **ignored**

### Profile example

Profile location: `./profile/jellyterra.json`

```json
{
  "recording": true,
  "webhooks": {
    "beforeConnect": [
      {
        "method": "POST",
        "url": "https://api.live.bilibili.com/room/v1/Room/startLive",
        "headers": {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        "cookies": {
          "bili_jct": "123456",
          "SESSDATA": "ABCDEF="
        },
        "body": "room_id=&area_v2=&platform=pc_link&csrf=123456",
        "body_encoding": "string"
      }
    ]
  },
  "remotes": [
    {
      "host": "owncast.jellyterra.com",
      "path": "live",
      "key": "password"
    },
    {
      "host": "live-push.bilivideo.com",
      "path": "live-bvc/",
      "key": "?streamname=&key=&schedule=rtmp&pflag=1"
    },
    {
      "host": "ingest.global-contribute.live-video.net",
      "path": "app/",
      "key": "live_123456_twitch"
    }
  ]
}
```

`body_encoding` claims the format of string in `body`:

- `string`
- `base64`

### Example

Server: `rtmp://rtmp-pioneer.jellyterra.com/profile/jellyterra`

```
$ rtmp-pioneer -o ~/Videos
Listen on :1935
1749123456123456 Inbound - profile: jellyterra
1749123456123456 Webhook beforeConnect - requesting index 0 url: https://api.live.bilibili.com/room/v1/Room/startLive
1749123456123456 Webhook response 200 text: {"code":0,"data":{"status":"LIVE"}}
1749123456123456 Endpoint 0 connecting: owncast.jellyterra.com
1749123456123456 Endpoint 1 connecting: live-push.bilivideo.com
1749123456123456 Endpoint 2 connecting: ingest.global-contribute.live-video.net
1749123456123456 Streaming started.
1749123456123456 Closed.
```

# Webhook

RTMP Pioneer supports running webhooks on following events:

| Event           | Usage                                                |
|-----------------|------------------------------------------------------|
| `beforeConnect` | Notify your platform to accept streaming.            |
| `afterStart`    | After all endpoints connected and streaming started. | 
| `afterClose`    | After client connection closed.                      |
