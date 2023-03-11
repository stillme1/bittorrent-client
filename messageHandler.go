package main

import (
	"encoding/binary"
	"io"
	"net"
	"time"
)

func messageType(peerConnection *PeerConnection, t int) (int32, int32, error) {

	peerConnection.conn.SetDeadline(time.Now().Add(time.Duration(t) * time.Second))
	defer peerConnection.conn.SetDeadline(time.Time{})

	buff1 := make([]byte, 4)
	_, err := io.ReadFull(peerConnection.conn, buff1)
	if err != nil {
		// return nil error if timeout happens, else return error
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return 0, -2, nil
		}
		return -1, -1, err
	}

	msglen := int32(binary.BigEndian.Uint32(buff1))
	if msglen == 0 {
		return 0, -1, nil
	}

	buff2 := make([]byte, 1)
	_, err = io.ReadFull(peerConnection.conn, buff2)
	if err != nil {
		return -1, -1, err
	}

	msgId := int32(uint32(buff2[0]))

	return msglen, msgId, nil
}

func handleHave(peerConnection *PeerConnection, length int32) error {
	peerConnection.conn.SetDeadline(time.Now().Add(100 * time.Second))
	defer peerConnection.conn.SetDeadline(time.Time{})

	buff := make([]byte, length)
	_, err := io.ReadFull(peerConnection.conn, buff)
	if err != nil {
		return err
	}
	index := int32(binary.BigEndian.Uint32(buff))
	peerConnection.bitfield[index] = true

	return nil
}

func handleBitfield(peerConnection *PeerConnection, length int32) error {
	peerConnection.conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer peerConnection.conn.SetDeadline(time.Time{})

	buff := make([]byte, length)
	_, err := io.ReadFull(peerConnection.conn, buff)
	if err != nil {
		return err
	}
	for i, j := range buff {
		for bit := 0; bit < 8; bit++ {
			if (j&(1<<bit) != 0) && ((i+1)*8-bit-1 < len(peerConnection.bitfield)) {
				peerConnection.bitfield[(i+1)*8-bit-1] = true
			}
		}
	}
	return nil
}

func handleCancel(peerConnection *PeerConnection) {
	// TODO
}
func handlePort(peerConnection *PeerConnection) {
	// TODO
}
func handleRequest(peerConnection *PeerConnection) {
	// TODO
}
func handlePiece(peerConnection *PeerConnection, length int, piece []*Piece) error {
	peerConnection.conn.SetDeadline(time.Now().Add(50 * time.Second))
	defer peerConnection.conn.SetDeadline(time.Time{})

	buff := make([]byte, length)
	_, err := io.ReadFull(peerConnection.conn, buff)
	if err != nil {
		return err
	}
	ind := int32(binary.BigEndian.Uint32(buff[0:4]))
	offset := int32(binary.BigEndian.Uint32(buff[4:8]))
	copy(piece[ind].data[offset:], buff[8:])
	return nil
}

func handleMessage(peerConnection *PeerConnection, msgId, msgLength int32, piece []*Piece) error {
	switch msgId {
	case -1:
		// keep alive
		return nil
	case 0:
		// choke
		peerConnection.choked = true
	case 1:
		// unchoke
		peerConnection.choked = false
	case 2:
		// interested
		peerConnection.interested = true
	case 3:
		// not interested
		peerConnection.interested = false
	case 4:
		// have
		return handleHave(peerConnection, msgLength-1)
	case 5:
		// bitfield
		return handleBitfield(peerConnection, msgLength-1)
	case 6:
		// request
		// TODO
	case 7:
		// piece
		return handlePiece(peerConnection, int(msgLength-1), piece)
	case 8:
		// cancel
		// TODO
	case 9:
		// port
		// TODO
	default:
		println("Unknown message id: ", msgId, " with length: ", msgLength, "to the peer: ", peerConnection.peer.ip, ":", peerConnection.peer.port)
		peerConnection.conn.SetDeadline(time.Now().Add(100 * time.Second))
		defer peerConnection.conn.SetDeadline(time.Time{})
		buff := make([]byte, msgLength-1)
		_, err := io.ReadFull(peerConnection.conn, buff)
		return err
	}

	return nil
}
