package main

import (
	// Uncomment this line to pass the first stage
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
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
	decoded, err := altbencode.Decode(string(torrentData))
	baseMap := decoded.GetData().(map[string]altbencode.Node)

	info := baseMap["info"].GetData().(map[string]altbencode.Node)

	infoBencoded, err := altbencode.Encode(baseMap["info"]);
	infoBytes := []byte(infoBencoded)
	infoHash := sha1.Sum([]byte(infoBytes))

	announceUrl := baseMap["announce"].GetData().(string)


	peer_id := "-PC0001-123456789012"

	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, announceUrl, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	q := req.URL.Query()

	q.Add("info_hash", fmt.Sprintf("%s", infoHash))
	q.Add("peer_id", peer_id)
	q.Add("port", "6881")
	q.Add("uploaded", "0")
	q.Add("downloaded", "0")
	q.Add("left", fmt.Sprint(info["length"].GetData().(int)))
	q.Add("compact", "1")

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(string(responseBody))

	responseNodes, err := altbencode.Decode(string(responseBody))
	responseMap := responseNodes.GetData().(map[string]altbencode.Node)
	stringPeers := responseMap["peers"].GetData().(string)

	for i := 0; i < len(stringPeers); i += 6 {
		ip := stringPeers[i:i+4]
		port := stringPeers[i+4:i+6]
		fmt.Printf("%d.%d.%d.%d:%d\n", ip[0], ip[1], ip[2], ip[3], int(port[0])*256+ int(port[1]))
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
