package main

import (
	"encoding/binary"
	"encoding/hex"
)

func buildHandshake(infoHash string, peerId []byte) []byte {
	req := make([]byte, 68)

	// pstrlen
	copy(req[0:], []byte{19})

	// pstr
	copy(req[1:20], []byte("BitTorrent protocol"))

	// reserved
	copy(req[20:], []byte{0, 0, 0, 0, 0, 0, 0, 0})

	// info_hash
	info_hash,_ := hex.DecodeString(infoHash)
	copy(req[28:], info_hash)

	// peer_id
	copy(req[48:], peerId)
	
	return req
}

func sendKeepAlive(peerConn *PeerConnection) error {
	req := make([]byte, 4)

	// length
	copy(req[0:], []byte{0, 0, 0, 0})

	_, err := peerConn.conn.Write(req)
	return err
}

func sendChoke(peerConn *PeerConnection) error {
	req := make([]byte, 5)

	// length
	copy(req[0:], []byte{0, 0, 0, 1})

	// id
	copy(req[4:], []byte{0})

	_, err := peerConn.conn.Write(req)
	return err
}

func sendUnchoke(peerConn *PeerConnection) error {
	req := make([]byte, 5)

	// length
	copy(req[0:], []byte{0, 0, 0, 1})

	// id
	copy(req[4:], []byte{1})

	_, err := peerConn.conn.Write(req)
	return err
}

func sendInterested(peerConn *PeerConnection) error {
	req := make([]byte, 5)

	// length
	copy(req[0:], []byte{0, 0, 0, 1})

	// id
	copy(req[4:], []byte{2})

	_, err := peerConn.conn.Write(req)
	return err
}

func sendNotInterested(peerConn *PeerConnection) error {
	req := make([]byte, 5)

	// length
	copy(req[0:], []byte{0, 0, 0, 1})

	// id
	copy(req[4:], []byte{3})

	_, err := peerConn.conn.Write(req)
	return err
}

func sendHave(peerConn *PeerConnection, index uint32) error {
	req := make([]byte, 9)

	// length
	copy(req[0:], []byte{0, 0, 0, 5})

	// id
	copy(req[4:], []byte{4})

	// index
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	copy(req[5:], indexBytes)

	_, err := peerConn.conn.Write(req)
	return err
}

func sendBitfield(peerConn *PeerConnection, bitfield []byte) error {
	req := make([]byte, 5+len(bitfield))

	// length
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(bitfield)+1))
	copy(req[0:], length)

	// id
	copy(req[4:], []byte{5})

	// bitfield
	copy(req[5:], bitfield)

	_, err := peerConn.conn.Write(req)
	return err
}

func sendRequest(peerConn *PeerConnection, index uint32, offset uint32, length uint32) error {
	req := make([]byte, 17)

	// length
	copy(req[0:], []byte{0, 0, 0, 13})

	// id
	copy(req[4:], []byte{6})

	// index
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	copy(req[5:], indexBytes)

	// offset
	offsetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(offsetBytes, offset)
	copy(req[9:], offsetBytes)

	// length
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, length)
	copy(req[13:], lengthBytes)
	
	_, err := peerConn.conn.Write(req)
	if(err == nil) {
		println("send request", index, offset, length)
	} else {
		println("error = ", err.Error())
	}
	return err
}

func sendPiece(peerConn *PeerConnection, index uint32, offset uint32, block []byte) error {
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

	// offset
	offsetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(offsetBytes, offset)
	copy(req[9:], offsetBytes)

	// block
	copy(req[13:], block)
	
	_, err := peerConn.conn.Write(req)
	return err
}

func sendCancel(peerConn *PeerConnection, index uint32, offset uint32, length uint32) error {
	req := make([]byte, 17)

	// length
	copy(req[0:], []byte{0, 0, 0, 13})

	// id
	copy(req[4:], []byte{8})

	// index
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	copy(req[5:], indexBytes)

	// offset
	offsetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(offsetBytes, offset)
	copy(req[9:], offsetBytes)

	// length
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, length)
	copy(req[13:], lengthBytes)

	_, err := peerConn.conn.Write(req)
	return err
}

func sendPort(peerConn *PeerConnection, port uint16) error {
	req := make([]byte, 7)

	// length
	copy(req[0:], []byte{0, 0, 0, 3})

	// id
	copy(req[4:], []byte{9})

	// port
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, port)
	copy(req[5:], portBytes)

	_, err := peerConn.conn.Write(req)
	return err
}