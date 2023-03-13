package main

import (
	"fmt"
	"sync"
)

// constants
var piecelength = int64(0)
var PEER_ID = make([]byte, 20)
var pieceDone = make(map[int]bool)
var listOfPeers = make(map[string]bool)
var mutex = &sync.Mutex{}
var pieces []*Piece
var path string
var info bencodeTorrent

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func removePeer(peer Peer) {
	mutex.Lock()
	delete(listOfPeers, peer.ip + fmt.Sprintf("%v", peer.port))
	mutex.Unlock()
}

func deletePiece(k int) {
	mutex.Lock()
	pieces[k].data = nil
	mutex.Unlock()
}