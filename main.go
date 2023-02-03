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
	torrent,err := gotorrentparser.Parse(file)
	if err != nil {
		panic(err)
	}

	// getting the peers from the UDP trackers
	peers := getPeer(torrent, PEER_ID)
	var PeerConnection []PeerConnection		// live connection for each active peer

	// handshaking with each peer
	for _, i := range peers {
		go handShake(torrent, i, PEER_ID , &PeerConnection);
	}
	time.Sleep(10*time.Second)

	// getting the bitfield of each peer
	for _, i := range PeerConnection {
		i.bitfield = make([]byte, len(pieces))
	}
}
