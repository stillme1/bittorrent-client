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

func handShake(torrent *gotorrentparser.Torrent, peer Peer, pieces *[]Piece, workQueue chan *Piece) bool {
	conn, err := net.DialTimeout("tcp", peer.ip+":"+strconv.Itoa(int(peer.port)), 5*time.Second)
	if err != nil {
		return false
	}

	conn.Write(buildHandshake(torrent.InfoHash, PEER_ID))

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})
	resp := make([]byte, 68)
	_, err = io.ReadFull(conn, resp)

	if err != nil {
		return false
	}
	bitfield := make([]bool, len(*pieces))
	peerConnection := PeerConnection{conn, peer, resp[48:], true, false, &bitfield}
	go startDownload(&peerConnection, torrent, pieces, workQueue)
	activePeers++
	return true
}

func handleAllPendingMessages(peerConnection *PeerConnection, piece *[]Piece, t int) bool {
	for {
		msgLength, msgId, err := messageType(peerConnection, t)
		if msgId == -2 {
			return true
		}
		if err != nil {
			return false
		}
		err = handleMessage(peerConnection, msgId, msgLength, piece)
		if err != nil {
			return false
		}
	}
}

func getPiece(peerConnection *PeerConnection, piece *[]Piece, k uint32) bool {
	block := ((*piece)[k].length + 0x00004000 - 1) / 0x00004000
	for block > 0 {
		msgLength, msgId, err := messageType(peerConnection, 1200)
		if err != nil || msgLength == -1 {
			return false
		}
		if msgId == 7 {
			err := handleMessage(peerConnection, msgId, msgLength, piece)
			if err != nil {
				return false
			}
			block--
		} else {
			err := handleMessage(peerConnection, msgId, msgLength, piece)
			if err != nil {
				return false
			}
		}
	}
	return validatePiece(&(*piece)[k])
}

func requestPiece(peerConnection *PeerConnection, piece *[]Piece, k uint32) bool {
	for i := 0; i < (*piece)[k].length; i += 0x00004000 {
		blockSize := min(0x00004000, uint32((*piece)[k].length-i))
		sendRequest(peerConnection, uint32((*piece)[k].index), uint32(i), blockSize)
	}
	valid := getPiece(peerConnection, piece, k)
	return valid
}

func validatePiece(piece *Piece) bool {
	res := sha1.Sum(piece.data) == piece.hash
	if !res {
		println("invalid piece", piece.index)
	}
	return res
}

func startDownload(peerConnection *PeerConnection, torrent *gotorrentparser.Torrent, pieces *[]Piece, workQueue chan *Piece) {
	defer decrementActivePeers()
	sendUnchoke(peerConnection)
	sendInterested(peerConnection)
	handleAllPendingMessages(peerConnection, pieces, 5)

	for {
		for piece := range workQueue {
			if !(*peerConnection.bitfield)[piece.index] || peerConnection.choked {
				workQueue <- piece
				if peerConnection.choked {
					active := handleAllPendingMessages(peerConnection, pieces, 2)
					if !active {
						peerConnection.conn.Close()
						println("Connection closed first: ", peerConnection.peer.ip)
						rebuilt := rebuildHandShake(torrent, peerConnection.peer, peerConnection.peerId, peerConnection)
						if !rebuilt {
							return
						}
						println("Connection rebuilt: ", peerConnection.peer.ip)
					}
				}
				continue
			}
			// request multiple peers for last two pieces
			if len(workQueue) == 0 {
				workQueue <- piece
				for i := 0; i < len(*pieces); i++ {
					if !pieceDone[i] {
						workQueue <- &(*pieces)[i]
					}
				}
			}
			if pieceDone[piece.index] {
				<-workQueue
				continue
			}
			println("Requesting piece: " + strconv.Itoa(piece.index))
			valid := requestPiece(peerConnection, pieces, uint32(piece.index))
			if valid {
				mutex.Lock()
				pieceDone[piece.index] = true
				mutex.Unlock()
				println("recieved piece: ", piece.index, " ", len(pieceDone))
				sendHave(peerConnection, uint32(piece.index))
			} else {
				workQueue <- piece
				peerConnection.conn.Close()
				println("Connection closed second:", peerConnection.peer.ip)
				rebuilt := rebuildHandShake(torrent, peerConnection.peer, peerConnection.peerId, peerConnection)
				if !rebuilt {
					return
				}
				println("Rebuilt connection", peerConnection.peer.ip)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
