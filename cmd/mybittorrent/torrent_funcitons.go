package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/altbencode"
)

type TorrentDetails struct {
	Announce string 
	Length int 
	InfoHash [20]byte
	PieceLength int 
	PieceHashes [][]byte 
	Peers []TorrentPeer 
}

type TorrentPeer struct {
	IP string 
	Port int 
}

type EncodedHandshake struct {
	Body [68]byte
}

type DecodedHandshake struct {
	InfoHash [20]byte
	PeerId [20]byte
}

func (d *DecodedHandshake) Encode() EncodedHandshake {
	var encoded EncodedHandshake
	copy(encoded.Body[0:20], []byte("\x13BitTorrent protocol"))
	copy(encoded.Body[28:48], d.InfoHash[:])
	copy(encoded.Body[48:68], d.PeerId[:])
	return encoded
}

func (e *EncodedHandshake) Decode() DecodedHandshake {
	var decoded DecodedHandshake
	copy(decoded.InfoHash[:], e.Body[28:48])
	copy(decoded.PeerId[:], e.Body[48:68])
	return decoded
}

func handshake(torrentDetails TorrentDetails, peerAddress string) (DecodedHandshake, error) {
	
	// peer := torrentDetails.Peers[peerIndex]
	// conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peer.IP, peer.Port))

	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		fmt.Println(err)
		return DecodedHandshake{}, err
	}

	defer conn.Close()

	var decodedConMessage DecodedHandshake

	decodedConMessage.InfoHash = torrentDetails.InfoHash
	copy(decodedConMessage.PeerId[:], []byte("-PC0001-123456789012"))

	encodedConMessage := decodedConMessage.Encode()

	_, err = conn.Write(encodedConMessage.Body[:])

	if err != nil {
		fmt.Println(err)
		return DecodedHandshake{}, err
	}

	var response EncodedHandshake

	_, err = conn.Read(response.Body[:])

	if err != nil {
		fmt.Println(err)
		return DecodedHandshake{}, err
	}

	return response.Decode(), nil
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
	torrentDetails.PieceHashes = pieces

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

	}

}

