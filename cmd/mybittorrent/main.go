package main

import (
	// Uncomment this line to pass the first stage
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/altbencode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)


func info(torrentData []byte) {
	decoded, err := altbencode.Decode(string(torrentData))
	baseMap := decoded.GetData().(map[string]altbencode.Node)

	if err != nil {
		fmt.Println(err)
		return
	}

	// print the announce key
	announce := baseMap["announce"]
	fmt.Println("Tracker URL: " +  announce.GetData().(string))

	info := baseMap["info"].GetData().(map[string]altbencode.Node)

	fmt.Println("Length: " + fmt.Sprint((info["length"]).GetData().(int)))

	infoBencoded, err := altbencode.Encode(baseMap["info"]);

	// avoid thorny issues with encoding glyphs
	infoBytes := []byte(infoBencoded)
	// use crypto/sha1 to hash the info bencoded string
	infoHash := sha1.Sum([]byte(infoBytes));

	
	fmt.Printf("Info Hash: %x\n", infoHash)
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
	} else if command == "info" {
		torrentFilePath := os.Args[2]
		torrentData, err := os.ReadFile(torrentFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		info(torrentData)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
