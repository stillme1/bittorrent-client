package main

import (
	"encoding/binary"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
)

func buildHandshake(torrent *gotorrentparser.Torrent, peerId []byte) []byte {
	req := make([]byte, 68)

	// pstrlen
	copy(req[0:], []byte{19})

	// pstr
	copy(req[1:20], []byte("BitTorrent protocol"))

	// reserved
	copy(req[20:], []byte{0, 0, 0, 0, 0, 0, 0, 0})

	// info_hash
	copy(req[28:], torrent.InfoHash)

	// peer_id
	copy(req[48:], peerId)

	return req
}

func buildKeepAlive() []byte {
	req := make([]byte, 4)

	// length
	copy(req[0:], []byte{0, 0, 0, 0})

	return req
}

func buildChoke() []byte {
	req := make([]byte, 5)

	// length
	copy(req[0:], []byte{0, 0, 0, 1})

	// id
	copy(req[4:], []byte{0})

	return req
}

func buildUnchoke() []byte {
	req := make([]byte, 5)

	// length
	copy(req[0:], []byte{0, 0, 0, 1})

	// id
	copy(req[4:], []byte{1})

	return req
}

func buildInterested() []byte {
	req := make([]byte, 5)

	// length
	copy(req[0:], []byte{0, 0, 0, 1})

	// id
	copy(req[4:], []byte{2})

	return req
}

func buildNotInterested() []byte {
	req := make([]byte, 5)

	// length
	copy(req[0:], []byte{0, 0, 0, 1})

	// id
	copy(req[4:], []byte{3})

	return req
}

func buildHave(index uint32) []byte {
	req := make([]byte, 9)

	// length
	copy(req[0:], []byte{0, 0, 0, 5})

	// id
	copy(req[4:], []byte{4})

	// index
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	copy(req[5:], indexBytes)

	return req
}

func buildBitfield(bitfield []byte) []byte {
	req := make([]byte, 5+len(bitfield))

	// length
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(bitfield)+1))
	copy(req[0:], length)

	// id
	copy(req[4:], []byte{5})

	// bitfield
	copy(req[5:], bitfield)

	return req
}

func buildRequest(index uint32, begin uint32, length uint32) []byte {
	req := make([]byte, 17)

	// length
	copy(req[0:], []byte{0, 0, 0, 13})

	// id
	copy(req[4:], []byte{6})

	// index
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	copy(req[5:], indexBytes)

	// begin
	beginBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(beginBytes, begin)
	copy(req[9:], beginBytes)

	// length
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, length)
	copy(req[13:], lengthBytes)

	return req
}

func buildPiece(index uint32, begin uint32, block []byte) []byte {
	req := make([]byte, 13+len(block))

	// length
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(block)+9))
	copy(req[0:], length)

	// id
	copy(req[4:], []byte{7})

	// index
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	copy(req[5:], indexBytes)

	// begin
	beginBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(beginBytes, begin)
	copy(req[9:], beginBytes)

	// block
	copy(req[13:], block)

	return req
}

func buildCancel(index uint32, begin uint32, length uint32) []byte {
	req := make([]byte, 17)

	// length
	copy(req[0:], []byte{0, 0, 0, 13})

	// id
	copy(req[4:], []byte{8})

	// index
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	copy(req[5:], indexBytes)

	// begin
	beginBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(beginBytes, begin)
	copy(req[9:], beginBytes)

	// length
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, length)
	copy(req[13:], lengthBytes)

	return req
}

func buildPort(port uint16) []byte {
	req := make([]byte, 7)

	// length
	copy(req[0:], []byte{0, 0, 0, 3})

	// id
	copy(req[4:], []byte{9})

	// port
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, port)
	copy(req[5:], portBytes)

	return req
}