package main

import (
	"net"
)

type PeerConnection struct {
	conn       net.Conn
	peer       Peer
	peerId     []byte
	choked     bool
	interested bool
	bitfield   []*bool
}

type Piece struct {
	index  int
	length int
	hash   [20]byte
	done   bool
	data   []byte
}

type file struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Files       []file `bencode:"files"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type connResp struct {
	action        uint32
	transactionId uint32
	connectionId  uint64
}

type Peer struct {
	ip   string
	port uint16
}

type annResp struct {
	action        uint32
	transactionId uint32
	seeders       uint32
	leechers      uint32
	peers         []Peer
}
