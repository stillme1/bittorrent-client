package main

import (
	"net"
)

type PeerConnection struct {
	conn net.Conn
	peer Peer
	peerId []byte
	chocked bool
	interested bool
	bitfield []byte
}