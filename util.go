package main

import (
	"sync"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
)

// constants
var piecelength = 0
var PEER_ID = make([]byte, 20)
var activePeers = 0
var pieceDone = make(map[int]bool)
var listOfPeers = make(map[string]bool)
var mutex = &sync.Mutex{}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func decrementActivePeers() {
	mutex.Lock()
	activePeers--
	mutex.Unlock()
}

func removePeer(peer string) {
	mutex.Lock()
	delete(listOfPeers, peer)
	mutex.Unlock()
}

func getSize(torrent *gotorrentparser.Torrent) int64 {
	files := torrent.Files
	var size int64
	for _, val := range files {
		size += val.Length
	}
	return size
}