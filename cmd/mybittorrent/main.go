package main

import (
	// Uncomment this line to pass the first stage

	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/altbencode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)



func cmdInfo(torrentData []byte) {
	torrentDetails, err := parseTorrentFile(torrentData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Tracker URL: " +  torrentDetails.Announce)
	fmt.Println("Length: " + fmt.Sprint(torrentDetails.Length))
	fmt.Printf("Info Hash: %x\n", torrentDetails.InfoHash)
	fmt.Println("Piece Length: " + fmt.Sprint(torrentDetails.PieceLength))
	fmt.Println("Pieces Hashes:")
	for _, piece := range torrentDetails.PieceHashes {
		fmt.Printf("%x\n", piece)
	}

}

func cmdFetchPeers(torrentData []byte) {
	torrentDetails, err := parseTorrentFile(torrentData)
	if err != nil {
		fmt.Println(err)
		return
	}
	addPeersToTorrentDetails(&torrentDetails)
	for _, peer := range torrentDetails.Peers {
		fmt.Printf("%s:%d\n", peer.IP, peer.Port)
	}
}

func cmdHandShake(torrentData []byte,  peerAddress string) {
	torrentDetails, err := parseTorrentFile(torrentData)
	if err != nil {
		fmt.Println(err)
		return
	}
	addPeersToTorrentDetails(&torrentDetails)
	decodedHandshake, _ := handshake(torrentDetails, peerAddress)
	fmt.Printf("Peer ID: %x\n", decodedHandshake.PeerId)
}

func cmdDownloadPiece(torrentData []byte, outputFilename string, pieceId string) {
	torrentDetails, err := parseTorrentFile(torrentData)
	if err != nil {
		fmt.Println(err)
		return
	}
	addPeersToTorrentDetails(&torrentDetails)
	pieceIndex, err := strconv.Atoi(pieceId)
	if err != nil {
		fmt.Println(err)
		return
	}
	// just use the first peer for now
	peer := torrentDetails.Peers[0]
	downloadPiece(&torrentDetails, peer, pieceIndex, outputFilename)
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		
		decoded, err := altbencode.Decode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else if command == "peers" {
		torrentFilePath := os.Args[2]
		torrentData, err := os.ReadFile(torrentFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		cmdFetchPeers(torrentData)

	} else if command == "info" {
		torrentFilePath := os.Args[2]
		torrentData, err := os.ReadFile(torrentFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		cmdInfo(torrentData)
	} else if command == "handshake" {
		peerAddress := os.Args[3]
		torrentFilePath := os.Args[2]
		torrentData, err := os.ReadFile(torrentFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		cmdHandShake(torrentData, peerAddress)
	} else if command == "download_piece" {
		// os.Args[2] is just "-o"
		outputFilename := os.Args[3]
		torrentFilePath := os.Args[4]
		pieceId := os.Args[5]
		torrentData, err := os.ReadFile(torrentFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		cmdDownloadPiece(torrentData, outputFilename, pieceId)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
