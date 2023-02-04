package main

import (
	"encoding/binary"
	"io"
	"time"
)

func messageType(peerConnection *PeerConnection) (int32, int32, error) {

	peerConnection.conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer peerConnection.conn.SetDeadline(time.Time{})

	buff1 := make([]byte, 4)
	_,err := io.ReadFull(peerConnection.conn, buff1)
	if(err != nil) {
		return -1, -1, err
	}

	buff2 := make([]byte, 1)
	_,err = io.ReadFull(peerConnection.conn, buff2)
	if(err != nil) {
		return -1, -1, err
	}
	return int32(binary.BigEndian.Uint32(buff1)), int32(uint32(buff2[0])), nil
}

func handleHave(peerConnection *PeerConnection, length int32) error {
	peerConnection.conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer peerConnection.conn.SetDeadline(time.Time{})

	buff := make([]byte, length)
	_,err := io.ReadFull(peerConnection.conn, buff)
	if(err != nil) {
		peerConnection.conn.Close()
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
	_,err := io.ReadFull(peerConnection.conn, buff)
	if(err != nil) {
		peerConnection.conn.Close()
		return err
	}
	for i, j := range buff {
		for bit := 0; bit < 8; bit++ {
			if((j & (1 << bit) != 0) && ((i+1)*8 - bit - 1 < len(peerConnection.bitfield))) {
				peerConnection.bitfield[(i+1)*8 - bit - 1] = true
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
func handlePiece(peerConnection *PeerConnection, length int) ([] byte , error) {
	peerConnection.conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer peerConnection.conn.SetDeadline(time.Time{})

	buff := make([]byte, length)
	_,err := io.ReadFull(peerConnection.conn, buff)
	if(err != nil) {
		return nil, err
	}
	return buff[8:], nil
}

func handlePeerConnection(peerConnection *PeerConnection, buff *[]byte) bool{
	msgLength, msgId, err := messageType(peerConnection)

	if(err == io.EOF || msgLength == -1 || msgId == -1) {
		return false
	}

	switch msgId {
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
		err = handleHave(peerConnection, msgLength-1)
		if(err != nil) {
			peerConnection.conn.Close()
			return false
		}
	case 5:
		// bitfield
		err = handleBitfield(peerConnection, msgLength-1)
		if(err != nil) {
			peerConnection.conn.Close()
			return false
		}
	case 6:
		// request
		// TODO
	case 7:
		// piece
		res, err := handlePiece(peerConnection, int(msgLength-1))
		if(err != nil) {
			buff = nil
			return false
		}
		*buff = res
		return true
	case 8:
		// cancel
		// TODO
	case 9:
		// port
		// TODO
	default:
		peerConnection.conn.SetDeadline(time.Now().Add(3 * time.Second))
		defer peerConnection.conn.SetDeadline(time.Time{})
		buff := make([]byte, msgLength-1)
		_,err := io.ReadFull(peerConnection.conn, buff)
		return err == nil

	}
	
	return true
}