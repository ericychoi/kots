#!/bin/sh

HOST=192.168.2.100:9091
MAGNETLINK=`go run kots.go -regex '^무한도전.+151226\.HDTV\.H264\.720p-WITH$' -show 무한도전`

transmission-remote "$HOST" -n transmission:transmission -a "$MAGNETLINK"
