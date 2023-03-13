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

func handShake(torrent *gotorrentparser.Torrent, peer Peer, workQueue chan *Piece) {
	conn, err := net.DialTimeout("tcp", peer.ip+":"+strconv.Itoa(int(peer.port)), 5*time.Second)
	if err != nil {
		removePeer(peer)
		return
	}

	conn.Write(buildHandshake(torrent.InfoHash, PEER_ID))

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})
	resp := make([]byte, 68)
	_, err = io.ReadFull(conn, resp)

	if err != nil {
		removePeer(peer)
		return
	}
	bitfield := make([]bool, len(pieces))
	peerConnection := PeerConnection{conn, peer, resp[48:], true, false, &bitfield}
	go download(&peerConnection, torrent, workQueue)
}

func handleAllPendingMessages(peerConnection *PeerConnection, t int) bool {
	for {
		msgLength, msgId, err := messageType(peerConnection, t)
		if msgId == -2 {
			return true
		}
		if err != nil {
			return false
		}
		err = handleMessage(peerConnection, msgId, msgLength)
		if err != nil {
			return false
		}
	}
}

func getPiece(peerConnection *PeerConnection, k int) bool {
	if(pieces[k].data == nil) {
		data := make([]byte, pieces[k].length)
		pieces[k].data = &data
	}
	block := (pieces[k].length + 0x00004000 - 1) / 0x00004000
	for block > 0 {
		msgLength, msgId, err := messageType(peerConnection, 1200)
		if err != nil {
			return false
		}
		if msgId == 7 {
			err := handleMessage(peerConnection, msgId, msgLength)
			if err != nil {
				return false
			}
			block--
		} else {
			err := handleMessage(peerConnection, msgId, msgLength)
			if err != nil {
				return false
			}
		}
	}
	return validatePiece(k)
}

func requestPiece(peerConnection *PeerConnection, k int) bool {
	for i := int64(0); i < pieces[k].length; i += 0x00004000 {
		blockSize := min(0x00004000, uint32(pieces[k].length-i))
		sendRequest(peerConnection, uint32(pieces[k].index), uint32(i), blockSize)
	}
	return getPiece(peerConnection, k)
}

func validatePiece(k int) bool {
	if(pieces[k].data == nil) {
		return false
	}
	res := sha1.Sum(*pieces[k].data) == pieces[k].hash
	if !res {
		println("invalid piece", pieces[k].index)
	}
	return res
}

func download(peerConnection *PeerConnection, torrent *gotorrentparser.Torrent, workQueue chan *Piece) {
	defer removePeer(peerConnection.peer)

	sendUnchoke(peerConnection)
	sendInterested(peerConnection)
	handleAllPendingMessages(peerConnection, 5)

	for piece := range workQueue {
		if !(*peerConnection.bitfield)[piece.index] || peerConnection.choked {
			workQueue <- piece
			if peerConnection.choked {
				active := handleAllPendingMessages(peerConnection, 2)
				if !active {
					peerConnection.conn.Close()
					println("Connection closed: ", peerConnection.peer.ip)
					rebuilt := rebuildHandShake(torrent, peerConnection.peer, peerConnection.peerId, peerConnection)
					if !rebuilt {
						return
					}
					println("Connection rebuilt: ", peerConnection.peer.ip)
				}
			}
			continue
		}
		mutex.Lock()
		if len(workQueue) == 0 {
			for i:= range pieces {
				if !pieceDone[i] {
					workQueue <- pieces[i]
				}
			}
		}
		if pieceDone[piece.index] {
			mutex.Unlock()
			continue;
		}
		mutex.Unlock()
		println("Requesting piece: " + strconv.Itoa(piece.index))
		valid := requestPiece(peerConnection, piece.index)
		if valid {
			write(piece.index)
			println("recieved piece: ", piece.index, " ", len(pieceDone))
			sendHave(peerConnection, uint32(piece.index))
			piece.data = nil
		} else {
			workQueue <- piece
			peerConnection.conn.Close()
			println("Connection closed:", peerConnection.peer.ip)
			rebuilt := rebuildHandShake(torrent, peerConnection.peer, peerConnection.peerId, peerConnection)
			if !rebuilt {
				return
			}
			println("Rebuilt connection", peerConnection.peer.ip)
		}
	}
}


func startDownload(torrent *gotorrentparser.Torrent, workQueue chan *Piece) {
	for	{
		peerList := getPeers(torrent)
		for _, peer := range peerList {
			go handShake(torrent, peer, workQueue)
		}
		time.Sleep(60 * time.Second)
	}
}