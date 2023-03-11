package main

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func decrementActivePeers(activePeers *int) {
	*activePeers--
}
