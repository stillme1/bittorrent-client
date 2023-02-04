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

func requestPiece(peerConnection *PeerConnection, piece *Piece) []byte{
	for i := 0; i < piece.length; i += 300 {
		blockSize := min(300, uint32(piece.length - i))
		block := make([]byte, blockSize)
		sendRequest(peerConnection, uint32(piece.index), uint32(i), blockSize)
	}
}

func validatePiece(piece *Piece) bool {
	return sha1.Sum(piece.data) == piece.hash
}

func startDownload(peerConnection *PeerConnection, workQueue chan *Piece, finished chan *Piece) error {

	sendUnchoke(peerConnection)
	sendInterested(peerConnection)
	
	for piece := range workQueue {
		if !peerConnection.bitfield[piece.index] {
			workQueue <- piece
			continue
		}
		buff := requestPiece(peerConnection, piece)

	}
	return nil
}