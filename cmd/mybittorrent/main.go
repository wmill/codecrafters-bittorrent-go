package main

import (
	// Uncomment this line to pass the first stage
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/altbencode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
// func decodeBencode(bencodedString string) (interface{}, error) {
// 	if unicode.IsDigit(rune(bencodedString[0])) {
// 		var firstColonIndex int

// 		for i := 0; i < len(bencodedString); i++ {
// 			if bencodedString[i] == ':' {
// 				firstColonIndex = i
// 				break
// 			}
// 		}

// 		lengthStr := bencodedString[:firstColonIndex]

// 		length, err := strconv.Atoi(lengthStr)
// 		if err != nil {
// 			return "", err
// 		}

// 		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
// 	} else if rune(bencodedString[0])  == 'i' {
// 		length := len(bencodedString)
// 		if bencodedString[length-1] != 'e' {
// 			return "", fmt.Errorf("Invalid integer")
// 		}
// 		stringInt := bencodedString[1:length-1]
// 		integer, err := strconv.Atoi(stringInt)
// 		return integer, err
// 	} else {
// 		return "", fmt.Errorf("Only strings are supported at the moment")
// 	}
// }

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
	//fmt.Println("Length: " + fmt.Sprint(length.GetData().(int)))


	
	// jsonOutput, _ := json.Marshal(decoded)
	// fmt.Println(string(jsonOutput))

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
