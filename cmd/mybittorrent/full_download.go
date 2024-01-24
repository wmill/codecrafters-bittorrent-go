package main

import (
	"fmt"
	"os"
)

// This file manages the full download of a torrent file. It uses downloadPiece to download each piece of the file, and then writes the pieces to the output file.

func downloadFile(torrentDetails *TorrentDetails, peer *TorrentPeer, outputFilename string) {
	file, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for i := 0; i < len(torrentDetails.PieceHashes); i++ {
		piece, err := downloadPiece(torrentDetails, peer, i)
		if err != nil {
			fmt.Println(err)
			return
		}

		file.Write(piece.Body)
	}
}