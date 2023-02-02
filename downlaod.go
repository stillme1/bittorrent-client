package main

import (
	"io"
	"net"
	"strconv"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
)

func handShake(torrent *gotorrentparser.Torrent, peer Peer, peedId []byte , peerConnection *[]PeerConnection) bool {
	conn, err := net.DialTimeout("tcp", peer.ip+":"+strconv.Itoa(int(peer.port)), 5*time.Second)
	if err != nil {
		return false;
	}
	
	conn.Write(buildHandshake(torrent.InfoHash, peedId))

	conn.SetDeadline(time.Now().Add(3*time.Second))
	resp := make([]byte, 68)
	_,err = io.ReadFull(conn, resp)
	conn.SetDeadline(time.Time{})

	if err != nil {
		return false;
	}
	
	*peerConnection = append(*peerConnection, PeerConnection{conn, peer, resp[48:], true, false, nil})
	return true;
}