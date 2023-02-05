package main

import (
	"crypto/sha1"
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
	defer conn.SetDeadline(time.Time{})
	resp := make([]byte, 68)
	_,err = io.ReadFull(conn, resp)

	if err != nil {
		return false;
	}
	
	*peerConnection = append(*peerConnection, PeerConnection{conn, peer, resp[48:], true, false, nil})
	return true;
}

func handleAllPendingMessages(peerConnection *PeerConnection, buff *[]byte) {
	curr := true
	for curr {
		curr = handlePeerConnection(peerConnection, buff)
	}
}

func requestPiece(peerConnection *PeerConnection, piece *Piece) bool {
	
	buff := make([]byte, piece.length)
	for i := 0; i < piece.length; i += 0x00004000 {
		blockSize := min(0x00004000, uint32(piece.length - i))
		block := make([]byte, blockSize)
		sendRequest(peerConnection, uint32(piece.index), uint32(i), blockSize)
		handleAllPendingMessages(peerConnection, &block)
		copy(buff[i:], block)
	}
	return validatePiece(piece)
}

func validatePiece(piece *Piece) bool {
	return sha1.Sum(piece.data) == piece.hash
}

func startDownload(peerConnection *PeerConnection, workQueue chan *Piece, finished chan *Piece) error {
	var temp []byte
	handleAllPendingMessages(peerConnection, &temp)
	sendUnchoke(peerConnection)
	sendInterested(peerConnection)
	
	for piece := range workQueue {
		if !peerConnection.bitfield[piece.index] {
			workQueue <- piece
			continue
		}
		println("Requesting piece: " + strconv.Itoa(piece.index))
		if(requestPiece(peerConnection, piece)) {
			finished <- piece
			println("recieved piece: ", piece.index, " ", len(finished))
		} else {
			workQueue <- piece
			println("failed to recieve piece: " + strconv.Itoa(piece.index))
		}
	}
	return nil
}