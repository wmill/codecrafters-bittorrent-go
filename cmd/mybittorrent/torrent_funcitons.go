package main

import (
	"crypto/sha1"
	"fmt"

	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/altbencode"
)

type TorrentDetails struct {
	Announce string 
	Length int 
	InfoHash [20]byte
	PieceLength int 
	Pieces [][]byte 
	Peers []TorrentPeer 
}

type TorrentPeer struct {
	IP string 
	Port int 
}

func parseTorrentFile(torrentData []byte) (TorrentDetails, error) {
	var torrentDetails TorrentDetails

	decoded, err := altbencode.Decode(string(torrentData))
	baseMap := decoded.GetData().(map[string]altbencode.Node)
	info := baseMap["info"].GetData().(map[string]altbencode.Node)

	if err != nil {
		fmt.Println(err)
		return torrentDetails, err
	}

	infoBencoded, err := altbencode.Encode(baseMap["info"]);

	// avoid thorny issues with encoding glyphs
	infoBytes := []byte(infoBencoded)
	// use crypto/sha1 to hash the info bencoded string
	infoHash := sha1.Sum([]byte(infoBytes))

	pieceLength := info["piece length"].GetData().(int)
	piecesString := info["pieces"].GetData().(string)
	var pieces [][]byte
	for i := 0; i < len(piecesString); i += 20 {
		// var piece []byte
		// copy(piece[:], piecesString[i:i+20])
		pieces = append(pieces, []byte(piecesString[i:i+20]))
	}

	torrentDetails.Announce = baseMap["announce"].GetData().(string)
	torrentDetails.Length = info["length"].GetData().(int)
	torrentDetails.InfoHash = infoHash
	torrentDetails.PieceLength = pieceLength
	torrentDetails.Pieces = pieces

	return torrentDetails, nil
}