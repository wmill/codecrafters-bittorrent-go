package altbencode

import (
	"encoding/json"
	"sort"
)



type MapNode struct {
	Children map[string]Node
}

func (n MapNode) GetData() interface{} {
	return n.Children
}

//add a MarshalJSON method to MapNode
func (n MapNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Children)
}

func decodeDictionary(bencodedString string, startIndex int) (ParseResult, error) {
	var children map[string]Node
	children = make(map[string]Node)

	index := startIndex + 1
	for index < len(bencodedString) {
		if bencodedString[index] == 'e' {
			break
		}
		keyResult, err := decodeString(bencodedString, index)
		if err != nil {
			return ParseResult{}, err
		}
		key := keyResult.Node.(StringNode).Data
		index = keyResult.RemainingStringIndex
		valueResult, err := decode(bencodedString, index)
		if err != nil {
			return ParseResult{}, err
		}
		children[key] = valueResult.Node
		index = valueResult.RemainingStringIndex
	}
	return ParseResult{
		Node: MapNode{
			Children: children,
		},
		RemainingStringIndex: index + 1,
	}, nil
}

func encodeDictionary(dictionary Node) (string, error) {
	var result string
	result += "d"

	theMap := dictionary.GetData().(map[string]Node)
	keys := make([]string, 0, len(theMap))
	for k := range theMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	for _, key := range keys {
		encodedKey, err := encodedString(StringNode{key})
		if err != nil {
			return "", err
		}
		encodedNode, err := encode(theMap[key])
		if err != nil {
			return "", err
		}
		result += encodedKey + encodedNode
	}
	result += "e"
	return result, nil
}