package main

import "sync"

// constants
var piecelength = 0
var PEER_ID = make([]byte, 20)
var activePeers = 0
var pieceDone map[int]bool
var mutex = &sync.Mutex{}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func decrementActivePeers() {
	activePeers--
}
