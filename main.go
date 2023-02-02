package main

import (
	"math/rand"
	"os"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
)

func main() {
	PEER_ID := make([]byte, 20)
	rand.Read(PEER_ID)

	arg := os.Args[1:]
	torrent,err := gotorrentparser.ParseFromFile(arg[0])

	if err != nil {
		panic(err)
	}
	
	peers := getPeer(torrent, PEER_ID)
	var PeerConnection []PeerConnection

	for _, i := range peers {
		go handShake(torrent, i, PEER_ID , &PeerConnection);
	}
	time.Sleep(10*time.Second)
	println(len(PeerConnection))
}