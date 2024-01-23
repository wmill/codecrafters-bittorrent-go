package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"os"
)

const (
	Choke byte = 0
	Unchoke  = 1
	Interested  = 2
	NotInterested  = 3
	Have  = 4
	Bitfield  = 5
	Request  = 6
	Piece  = 7
	Cancel  = 8
)

const ChunkSize = 16384

type DecodedPeerMessage struct { 
	ID uint8
	Payload []byte
}

type EncodedPeerMessage struct {
	Body []byte
}

func (d *DecodedPeerMessage) Encode() EncodedPeerMessage {
	var encoded EncodedPeerMessage
	encoded.Body = make([]byte, len(d.Payload) + 5)
	length := uint32(len(d.Payload) + 1)
	lengthArray := make([]byte, 4)


	binary.BigEndian.PutUint32(lengthArray, length)

	encoded.Body[4] = d.ID
	copy(encoded.Body[0:4], lengthArray)
	copy(encoded.Body[5:], d.Payload)

	return encoded
}

func (e *EncodedPeerMessage) Decode() DecodedPeerMessage {
	var decoded DecodedPeerMessage
	length := binary.BigEndian.Uint32(e.Body[0:4])
	if (length + 4 != uint32(len(e.Body))) {
		fmt.Println("Invalid message length")
	}
	decoded.ID = e.Body[4]
	decoded.Payload = e.Body[5:]
	return decoded
}


func sendHandshake(torrentDetails *TorrentDetails, peerId string, conn net.Conn) (DecodedHandshake, error)  {
	// send handshake
	var decodedHandshake DecodedHandshake

	decodedHandshake.InfoHash = torrentDetails.InfoHash
	copy(decodedHandshake.PeerId[:], []byte(peerId))

	encodedHandshake := decodedHandshake.Encode()

	_, err := conn.Write(encodedHandshake.Body[:])

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

func recievePeerMessage(conn net.Conn) (DecodedPeerMessage, error) {
	var response EncodedPeerMessage


	for {
		lengthBuffer := make([]byte, 4)
		// _, err := conn.Read(lengthBuffer)
		_, err := io.ReadFull(conn, lengthBuffer)

		if err != nil {
			fmt.Println(err)
			if err.Error() == "EOF" {
				return DecodedPeerMessage{}, err
			}
			continue
		}

		length := binary.BigEndian.Uint32(lengthBuffer)
		// fmt.Printf("Length: %d\n", length)
		if length == 0 {
			fmt.Println("Keep alive")
			continue
		}

		response.Body = make([]byte, length + 4)
		copy(response.Body[0:4], lengthBuffer)
		// _, err = conn.Read(response.Body[4:])
		_, err = io.ReadFull(conn, response.Body[4:])
		// fmt.Printf("Read: %d\n", n)
		// fmt.Printf("Read Length: %d\n", len(response.Body[4:]))

		if err != nil {
			fmt.Println(err)
			if err.Error() == "EOF" {
				return DecodedPeerMessage{}, err
			}
			continue
		}

		break

	}

	return response.Decode(), nil
}


func downloadPiece(torrentDetails *TorrentDetails, peer TorrentPeer, pieceIndex int, outputFilename string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peer.IP, peer.Port))
	// fmt.Printf("%s:%d", peer.IP, peer.Port)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	peerId := "-PC0001-123456789012"

	sendHandshake(torrentDetails, peerId, conn)

	response, err := recievePeerMessage(conn)

	if response.ID != Bitfield {
		fmt.Println("Expected bitfield message")
		fmt.Println(response.ID)
	}


	// send interested message
	var decodedInterestedMessage DecodedPeerMessage
	decodedInterestedMessage.ID = Interested
	encodedInterestedMessage := decodedInterestedMessage.Encode()
	_, err = conn.Write(encodedInterestedMessage.Body[:])

	if err != nil {
		fmt.Println(err)
		return
	}

	// receive unchoke message
	
	response, err = recievePeerMessage(conn)

	if response.ID != Unchoke {
		fmt.Println("Expected unchoke message")
		fmt.Println(response.ID)
		return
	}


	pieceLength := torrentDetails.PieceLength
	numberOfPieces := len(torrentDetails.PieceHashes)
	if pieceIndex == numberOfPieces - 1 {
		pieceLength = torrentDetails.Length - (numberOfPieces - 1) * torrentDetails.PieceLength
	}

	numberOfChunks := int(math.Ceil(float64(pieceLength) / float64(ChunkSize)))
	// numberOfChunks := torrentDetails.PieceLength / ChunkSize + 1
	pieceData := ""

	retryCount := 5
	
	for i := 0; i < numberOfChunks; i++ {
		// send request message
		targetChunkSize :=  ChunkSize
		if i == numberOfChunks - 1 {
			targetChunkSize = pieceLength - (numberOfChunks - 1) * ChunkSize
		}
		if targetChunkSize == 0 {
			break
		}
		var decodedRequestMessage DecodedPeerMessage
		decodedRequestMessage.ID = Request
		decodedRequestMessage.Payload = make([]byte, 12)
		binary.BigEndian.PutUint32(decodedRequestMessage.Payload[0:4], uint32(pieceIndex))
		binary.BigEndian.PutUint32(decodedRequestMessage.Payload[4:8], uint32(i * ChunkSize))
		binary.BigEndian.PutUint32(decodedRequestMessage.Payload[8:12], uint32(targetChunkSize))
		encodedRequestMessage := decodedRequestMessage.Encode()
		_, err = conn.Write(encodedRequestMessage.Body[:])

		if err != nil {
			fmt.Println(err)
			break
		}

		// receive piece message
		response, err = recievePeerMessage(conn)

		if err != nil {
			fmt.Println(err)
			if err.Error() == "EOF" {
				if (retryCount == 0) {
					break
				}
				retryCount--
				i--
				continue
			}
			break
		}

		// following responses don't have the piece id
		if response.ID == Piece {
			pieceData += string(response.Payload[8:])
			//copy(pieceData[i * ChunkSize:], response.Payload[8:])
		} else {
			// copy(pieceData[i * ChunkSize:], response.Payload)
			fmt.Println("Expected piece message, got something else")
			fmt.Println(response.ID)
			fmt.Printf("payload length: %d\n", len(response.Payload))
			// pieceData += string(response.Payload)
			// fmt.Println("Retyring")
			// i--
			// continue
		}
		// time.Sleep(100)
	
	}
	fmt.Printf("Piece size: %d\nData size: %d\n", pieceLength, len(pieceData))
	os.WriteFile(outputFilename, []byte(pieceData), 0644)
	// fmt.Printf("%d\n", len(pieceData))
	// fmt.Printf("%s\n", pieceData)
}