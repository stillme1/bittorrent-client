package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
	bencode "github.com/jackpal/bencode-go"
)

func main() {
	// Generating a random peer id
	rand.Read(PEER_ID)
	pieceDone = make(map[int]bool)

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
	pieces := make([]Piece, len(piecesString)/20)
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
		pieces[i/20] = temp
	}

	// partsing the torrent file using go-torrent-parser
	torrent, err := gotorrentparser.ParseFromFile(arg[0])
	if err != nil {
		panic(err)
	}

	workQueue := make(chan *Piece, len(pieces))

	for i := range pieces {
		workQueue <- &pieces[i]
	}

	// Starting download
	go startDownload(torrent, &pieces, workQueue)

	for len(pieceDone) != len(pieces) {
		fmt.Println("download = ", float64(len(pieceDone))/float64(len(pieces))*100, "%")
		fmt.Println("active peers = ", len(listOfPeers))
		time.Sleep(10 * time.Second)
	}

	if len(info.Info.Files) == 0 {
		singleFileWrite(info, &pieces, arg[1])
	} else {
		multiFileWrite(info, &pieces, arg[1])
	}
}
