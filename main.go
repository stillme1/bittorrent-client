package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
	bencode "github.com/jackpal/bencode-go"
)

var piecelength = 0
var PEER_ID = make([]byte, 20)

func main() {
	// Generating a random peer id
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
	piecelength = info.Info.PieceLength

	// getting the hash of each piece as hexadecimal string
	for i := range info.Info.Files {
		info.Info.Length += info.Info.Files[i].Length
	}

	lastpieceLength := info.Info.Length % info.Info.PieceLength
	piecesString := info.Info.Pieces
	pieces := make([]*Piece, len(piecesString)/20)
	for i := 0; i < len(piecesString); i += 20 {
		var temp Piece
		temp.index = i / 20
		for j := 0; j < 20; j++ {
			temp.hash[j] = piecesString[i+j]
		}
		if i+20 == len(piecesString) {
			temp.length = lastpieceLength
		} else {
			temp.length = info.Info.PieceLength
		}
		temp.data = make([]byte, temp.length)
		pieces[i/20] = &temp

	}

	// partsing the torrent file using go-torrent-parser
	torrent, err := gotorrentparser.ParseFromFile(arg[0])
	if err != nil {
		panic(err)
	}

	// getting the peers from the UDP trackers
	peers := getPeer(torrent, PEER_ID)
	var peerConnection []PeerConnection // live connection for each active peer

	// handshaking with each peer
	for _, i := range peers {
		go handShake(torrent, i, PEER_ID, &peerConnection)
	}
	time.Sleep(12 * time.Second)

	// getting the bitfield of each peer
	for i := range peerConnection {
		peerConnection[i].bitfield = make([]bool, len(pieces))
	}

	workQueue := make(chan *Piece, len(pieces))
	finishedQueue := make(chan *Piece, len(pieces))

	for i := range pieces {
		workQueue <- pieces[i]
	}
	activePeers := len(peerConnection)
	for i := range peerConnection {
		go startDownload(&peerConnection[i], torrent, pieces, &activePeers, workQueue, finishedQueue)
	}

	for len(finishedQueue) != len(pieces) {
		fmt.Println("download = ", float64(len(finishedQueue))/float64(len(pieces))*100, "%")
		fmt.Println("active peers = ", activePeers)
		time.Sleep(10 * time.Second)
	}

	if len(info.Info.Files) == 0 {
		singleFileWrite(info, pieces, arg[1])
	} else {
		multiFileWrite(info, pieces, arg[1])
	}
}
