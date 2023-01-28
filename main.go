package main

import (
	"math/rand"
	"os"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
)

func main() {
	PEER_ID := make([]byte, 20)
	rand.Read(PEER_ID)

	arg := os.Args[1:]
	println(arg[0])
	torrent,err := gotorrentparser.ParseFromFile(arg[0])

	if err != nil {
		panic(err)
	}
	
	peers := getPeer(torrent, PEER_ID)

	println("list of peers: ")
	for _, i := range peers {
		println(i.ip, i.port)
	}
}