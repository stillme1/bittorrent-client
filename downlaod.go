package main

import (
	"crypto/sha1"
	"io"
	"net"
	"strconv"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
)

func rebuildHandShake(torrent *gotorrentparser.Torrent, peer Peer, peedId []byte, peerConnection *PeerConnection) bool {
	conn, err := net.DialTimeout("tcp", peer.ip+":"+strconv.Itoa(int(peer.port)), 5*time.Second)
	if err != nil {
		return false
	}

	conn.Write(buildHandshake(torrent.InfoHash, peedId))

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})
	resp := make([]byte, 68)
	_, err = io.ReadFull(conn, resp)

	if err != nil {
		return false
	}
	peerConnection.conn = conn
	return true
}

func handShake(torrent *gotorrentparser.Torrent, peer Peer, peedId []byte, peerConnection *[]PeerConnection) bool {
	conn, err := net.DialTimeout("tcp", peer.ip+":"+strconv.Itoa(int(peer.port)), 5*time.Second)
	if err != nil {
		return false
	}

	conn.Write(buildHandshake(torrent.InfoHash, peedId))

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})
	resp := make([]byte, 68)
	_, err = io.ReadFull(conn, resp)

	if err != nil {
		return false
	}

	*peerConnection = append(*peerConnection, PeerConnection{conn, peer, resp[48:], true, false, nil})
	return true
}

func handleAllPendingMessages(peerConnection *PeerConnection, piece []*Piece, t uint32) bool {
	curr := true
	var err error
	for curr {
		curr, err = handlePeerConnection(peerConnection, piece, t)
	}
	return err == nil
}

func requestPiece(peerConnection *PeerConnection, piece []*Piece, k uint32) (bool, bool) {
	for i := 0; i < piece[k].length; i += 0x00004000 {
		blockSize := min(0x00004000, uint32(piece[k].length-i))
		sendRequest(peerConnection, uint32(piece[k].index), uint32(i), blockSize)
	}
	active := handleAllPendingMessages(peerConnection, piece, 10)
	return validatePiece(piece[k]), active
}

func validatePiece(piece *Piece) bool {
	res := sha1.Sum(piece.data) == piece.hash
	if !res {
		println("invalid piece", piece.index)
	}
	return res
}

func startDownload(peerConnection *PeerConnection, torrent *gotorrentparser.Torrent, pieces []*Piece, workQueue chan *Piece, finished chan *Piece) {
	sendUnchoke(peerConnection)
	sendInterested(peerConnection)
	handleAllPendingMessages(peerConnection, pieces, 1)

	for piece := range workQueue {
		if !peerConnection.bitfield[piece.index] || peerConnection.choked {
			workQueue <- piece
			if peerConnection.choked {
				active := handleAllPendingMessages(peerConnection, pieces, 1)
				if !active {
					println("Closing conn", peerConnection.peer.ip)
					peerConnection.conn.Close()
					rebuilt := rebuildHandShake(torrent, peerConnection.peer, peerConnection.peerId, peerConnection)
					if !rebuilt {
						return
					}
				}
			}
			continue
		}
		println("Requesting piece: " + strconv.Itoa(piece.index))
		valid, active := requestPiece(peerConnection, pieces, uint32(piece.index))
		if valid {
			finished <- piece
			println("recieved piece: ", piece.index, " ", len(finished))
			sendHave(peerConnection, uint32(piece.index))
		} else {
			workQueue <- piece
		}
		if !active {
			println("Closing conn", peerConnection.peer.ip)
			peerConnection.conn.Close()
			rebuilt := rebuildHandShake(torrent, peerConnection.peer, peerConnection.peerId, peerConnection)
			if !rebuilt {
				return
			}
		}
	}
	return
}
