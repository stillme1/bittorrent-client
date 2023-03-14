package main

import (
	"os"
	"path/filepath"
)

func write(k int) {
	mutex.Lock()
	defer mutex.Unlock()
	defer markPieceDone(k)
	defer deletePiece(k)
	if pieceDone[k] {
		return
	}
	offset := int64(k) * piecelength
	left := pieces[k].length
	i := 0
	for i < len(info.Info.Files) {
		if offset < info.Info.Files[i].Length {
			break
		}
		offset -= info.Info.Files[i].Length
		i++
	}
	dataOffset := int64(0)
	for i < len(info.Info.Files) {
		filePath := path + "/" + info.Info.Name
		for _, j := range info.Info.Files[i].Path {
			filePath += "/" + j
		}
		err := os.MkdirAll(filepath.Dir(filePath), 0777)
		if err != nil {
			panic(err)
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.Seek(offset, 0)
		if left <= info.Info.Files[i].Length - offset {
			_, err := file.Write((*pieces[k].data)[dataOffset : ])
			if err != nil {
				panic(err)
			}
			return
		}
		_, err = file.Write((*pieces[k].data)[dataOffset: dataOffset + info.Info.Files[i].Length - offset])
		if err != nil {
			panic(err)
		}
		dataOffset += info.Info.Files[i].Length - offset
		left -= info.Info.Files[i].Length - offset
		offset = 0
		i++
	}
}
