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

func getPiece(peerConnection *PeerConnection, piece []*Piece, k uint32) bool {
	block := (piece[k].length + 0x00004000 - 1) / 0x00004000
	for block > 0{
		msgLength, msgId, err := messageType(peerConnection, 500)
		if err != nil || msgLength == -1{
			return false
		}
		if msgId == 7 {
			active, err := handleMessage(peerConnection, msgId, msgLength, piece)
			if err != nil || !active {
				return false
			}
			block--
		} else {
			active, err := handleMessage(peerConnection, msgId, msgLength, piece)
			if err != nil || !active {
				return false
			}
		}
	}
	return validatePiece(piece[k])
}

func requestPiece(peerConnection *PeerConnection, piece []*Piece, k uint32) (bool) {
	for i := 0; i < piece[k].length; i += 0x00004000 {
		blockSize := min(0x00004000, uint32(piece[k].length-i))
		sendRequest(peerConnection, uint32(piece[k].index), uint32(i), blockSize)
	}
	active := getPiece(peerConnection, piece, k)
	return active
}

func validatePiece(piece *Piece) bool {
	res := sha1.Sum(piece.data) == piece.hash
	if !res {
		println("invalid piece", piece.index)
	}
	return res
}

func startDownload(peerConnection *PeerConnection, torrent *gotorrentparser.Torrent, pieces []*Piece, activePeers *int, workQueue chan *Piece, finished chan *Piece) {
	defer decrementActivePeers(activePeers)
	sendUnchoke(peerConnection)
	sendInterested(peerConnection)
	handleAllPendingMessages(peerConnection, pieces, 5)

	for {
		for piece := range workQueue {
			if !peerConnection.bitfield[piece.index] || peerConnection.choked {
				workQueue <- piece
				if peerConnection.choked {
					active := handleAllPendingMessages(peerConnection, pieces, 2)
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
			// request multiple peers for last two pieces
			if len(workQueue) < 2 {
				workQueue <- piece
			}
			println("Requesting piece: " + strconv.Itoa(piece.index))
			valid := requestPiece(peerConnection, pieces, uint32(piece.index))
			if valid {
				finished <- piece
				println("recieved piece: ", piece.index, " ", len(finished))
				sendHave(peerConnection, uint32(piece.index))
			} else {
				workQueue <- piece
				println("Closing conn", peerConnection.peer.ip)
				peerConnection.conn.Close()
				rebuilt := rebuildHandShake(torrent, peerConnection.peer, peerConnection.peerId, peerConnection)
				if !rebuilt {
					return
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}
