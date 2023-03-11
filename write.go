package main

import (
	"os"
	"path/filepath"
)

func singleFileWrite(info bencodeTorrent, pieces []*Piece, path string) {

	path += "/" + info.Info.Name
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		panic("Error creating file" + err.Error())
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic("Error creating file" + err.Error())
	}
	defer file.Close()

	offset := int64(0)
	for i := range pieces {
		file.Seek(offset, 0)
		_, err := file.Write(pieces[i].data)
		if err != nil {
			panic("Error writing file, " + err.Error())
		}
		offset += int64(pieces[i].length)
	}
}

func multiFileWrite(info bencodeTorrent, pieces []*Piece, path string) {

	currPiece := 0
	offset := int64(0)
	for _, i := range info.Info.Files {
		filePath := path + "/" + info.Info.Name
		for _, j := range i.Path {
			filePath += "/" + j
		}
		err := os.MkdirAll(filepath.Dir(filePath), 0777)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			panic("Error creating file" + err.Error())
		}
		foffset := int64(0)
		for i.Length > 0 {
			file.Seek(foffset, 0)
			k := int64(piecelength) - offset
			if k > int64(i.Length) {
				k = int64(i.Length)
			}
			file.Write(pieces[currPiece].data[offset : offset+k])
			i.Length -= int(k)
			offset += k
			if offset == int64(piecelength) {
				offset = 0
				currPiece++
			}
			foffset += k
		}
		file.Close()
	}
}
