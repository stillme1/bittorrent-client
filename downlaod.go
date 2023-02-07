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

func handleAllPendingMessages(peerConnection *PeerConnection, piece []*Piece) bool {
	curr := true
	var err error
	for curr {
		curr, err = handlePeerConnection(peerConnection, piece)
	}
	return err == nil
}

func requestPiece(peerConnection *PeerConnection, piece []*Piece, k uint32) bool {
	active := true
	for i := 0; i < piece[k].length && active; i += 0x00004000 {
		blockSize := min(0x00004000, uint32(piece[k].length - i))
		sendRequest(peerConnection, uint32(piece[k].index), uint32(i), blockSize)
		active = handleAllPendingMessages(peerConnection, piece)
	}
	return validatePiece(piece[k]) && active
}

func validatePiece(piece *Piece) bool {
	res := sha1.Sum(piece.data) == piece.hash
	if(!res) {
		println("invalid piece", piece.index)
	}
	return res
}

func startDownload(peerConnection *PeerConnection, pieces []*Piece, workQueue chan *Piece, finished chan *Piece) {
	handleAllPendingMessages(peerConnection, pieces)
	sendUnchoke(peerConnection)
	sendInterested(peerConnection)
	
	for piece := range workQueue {
		err := sendKeepAlive(peerConnection)
		if(err != nil) {
			return
		}
		if !peerConnection.bitfield[piece.index] || peerConnection.choked {
			if(peerConnection.choked) {
				handleAllPendingMessages(peerConnection, pieces)
			}
			workQueue <- piece
			continue
		}
		println("Requesting piece: " + strconv.Itoa(piece.index))
		if(requestPiece(peerConnection, pieces, uint32(piece.index))) {
			finished <- piece
			println("recieved piece: ", piece.index, " ", len(finished))
			sendHave(peerConnection, uint32(piece.index))
		} else {
			workQueue <- piece
		}
	}
	return
}