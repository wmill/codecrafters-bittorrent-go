package altbencode

import (
	"encoding/json"
	"strconv"
)

type StringNode struct {
	Data string
}

func (n StringNode) GetData() interface{} {
	return n.Data
}

//add a MarshalJSON method 
func (n StringNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Data)
}

func decodeString(bencodedString string, startIndex int) (ParseResult, error) {
	var firstColonIndex int

	for i := startIndex; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[startIndex:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return ParseResult{}, err
	}

	return ParseResult{
		Node: StringNode{
			Data: bencodedString[firstColonIndex+1 : firstColonIndex+1+length],
		},
		RemainingStringIndex: firstColonIndex+1+length,
	}, nil
}

func encodedString(str Node) (string, error) {
	return strconv.Itoa(len(str.GetData().(string))) + ":" + str.GetData().(string), nil
}
