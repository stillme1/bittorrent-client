package main

import (
	"encoding/hex"
	"math/rand"
	"os"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
	bencode "github.com/jackpal/bencode-go"
)

func main() {
	// Generating a random peer id
	PEER_ID := make([]byte, 20)
	rand.Read(PEER_ID)

	// getting the path to torrent file as an argument
	arg := os.Args[1:]
	file, err := os.Open(arg[0])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// parsing the torrent file
	info := bencodeTorrent{}
	err = bencode.Unmarshal(file, &info)
	if err != nil {
		panic(err)
	}

	// getting the hash of each piece as hexadecimal string
	var pieces []string
	piecesString := info.Info.Pieces
	for i := 0; i < len(piecesString); i += 20 {
		pieces = append(pieces, hex.EncodeToString([]byte(piecesString[i:i+20])))
	}

	// partsing the torrent file using go-torrent-parser
	torrent,err := gotorrentparser.ParseFromFile(arg[0])
	if err != nil {
		panic(err)
	}

	// getting the peers from the UDP trackers
	peers := getPeer(torrent, PEER_ID)
	var peerConnection []PeerConnection		// live connection for each active peer

	// handshaking with each peer
	for _, i := range peers {
		go handShake(torrent, i, PEER_ID , &peerConnection);
	}
	time.Sleep(12*time.Second)

	// getting the bitfield of each peer
	for i,_ := range peerConnection {
		peerConnection[i].bitfield = make([]bool, len(pieces))
	}

	// array to store current state of pieces
	// 0 -> not started
	// 1 -> completed
	// 2 -> in progress
	status := make([]int, len(pieces))

	for i,_ := range peerConnection {
		go startDownload(&peerConnection[i], &status)
		println(peerConnection[i].peer.ip)
	}
	time.Sleep(10 * time.Second)
}
