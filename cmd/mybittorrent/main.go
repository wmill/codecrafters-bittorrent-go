package main

import (
	// Uncomment this line to pass the first stage
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net"
	"os"

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
	for _, piece := range torrentDetails.Pieces {
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

func cmdHandShake(torrentData []byte, peerAddress string) {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	var conMessage []byte

	conMessage = append(conMessage, []byte("\x13BitTorrent protocol\x00\x00\x00\x00\x00\x00\x00\x00")...)

	decoded, err := altbencode.Decode(string(torrentData))
	baseMap := decoded.GetData().(map[string]altbencode.Node)


	infoBencoded, err := altbencode.Encode(baseMap["info"]);
	infoBytes := []byte(infoBencoded)
	infoHash := sha1.Sum([]byte(infoBytes))
	conMessage = append(conMessage, infoHash[:]...)
	conMessage = append(conMessage, []byte("00112233445566778899")...)

	conn.Write(conMessage)

	reply := make([]byte, 68)
	conn.Read(reply)
	// fmt.Println(reply)
	responsePeerId := reply[48:68]
	fmt.Printf("Peer ID: %x\n", responsePeerId)
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
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
