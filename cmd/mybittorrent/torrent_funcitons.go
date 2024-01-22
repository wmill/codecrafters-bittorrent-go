package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"

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
		pieces = append(pieces, []byte(piecesString[i:i+20]))
	}

	torrentDetails.Announce = baseMap["announce"].GetData().(string)
	torrentDetails.Length = info["length"].GetData().(int)
	torrentDetails.InfoHash = infoHash
	torrentDetails.PieceLength = pieceLength
	torrentDetails.Pieces = pieces

	return torrentDetails, nil
}

func addPeersToTorrentDetails(torrentDetails *TorrentDetails)  {
	peer_id := "-PC0001-123456789012"

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, torrentDetails.Announce, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	q := req.URL.Query()

	q.Add("info_hash", fmt.Sprintf("%s", torrentDetails.InfoHash))
	q.Add("peer_id", peer_id)
	q.Add("port", "6881")
	q.Add("uploaded", "0")
	q.Add("downloaded", "0")
	q.Add("left", fmt.Sprint(torrentDetails.Length))
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

	responseNodes, err := altbencode.Decode(string(responseBody))
	responseMap := responseNodes.GetData().(map[string]altbencode.Node)
	stringPeers := responseMap["peers"].GetData().(string)

	for i := 0; i < len(stringPeers); i += 6 {
		ip := stringPeers[i:i+4]
		port := stringPeers[i+4:i+6]
		stringIp := fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
		intPort := int(port[0])*256+ int(port[1])
		torrentDetails.Peers = append(torrentDetails.Peers, TorrentPeer{IP: stringIp, Port: intPort})
		//fmt.Printf("%d.%d.%d.%d:%d\n", ip[0], ip[1], ip[2], ip[3], int(port[0])*256+ int(port[1]))
	}

}