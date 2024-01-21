package bencode

import "strconv"

type Node interface {
	GetData() interface{}
}

type ParseResult struct {
	Node Node
	RemainingString string
}

func (r ParseResult) GetData() interface{} {
	return r.Node.GetData()
}

type StringNode struct {
	Data string
}

func (n StringNode) GetData() interface{} {
	return n.Data
}

type IntNode struct {
	Data int
}

func (n IntNode) GetData() interface{} {
	return n.Data
}

type ListNode struct {
	Children []Node
}

func (n ListNode) GetData() interface{} {
	return n.Children
}

type MapNode struct {
	Children map[string]Node
}

type InvalidBencodeError struct {
	Message string
}

func (e *InvalidBencodeError) Error() string {
	return e.Message
}

func (n MapNode) GetData() interface{} {
	return n.Children
}

func Decode(bencodedString string) (Node, error) {
	if len(bencodedString) == 0 {
		return nil, nil
	}

	switch bencodedString[0] {
	case 'i':
		return decodeInt(bencodedString)
	case 'l':
		return decodeList(bencodedString)
	case 'd':
		return decodeMap(bencodedString)
	default:
		return decodeString(bencodedString)
	}
}

func decodeInt(bencodedString string) (Node, error) {
	length := len(bencodedString)
	if bencodedString[length-1] != 'e' {
		return nil, &InvalidBencodeError{"Invalid integer"}
	}
	stringInt := bencodedString[1:length-1]
	integer, err := strconv.Atoi(stringInt)
	if err != nil {
		return nil, err
	}
	return IntNode{integer}, nil
}

func decodeList(bencodedString string) (Node, error) {
	listNode := ListNode{}
	remainingString := bencodedString[1:]

	for remainingString[0] != 'e' {
		node, err := Decode(remainingString)
		if err != nil {
			return nil, err
		}
		listNode.Children = append(listNode.Children, node)
		remainingString = node.(ParseResult).RemainingString
	}

	return listNode, nil
}	

func decodeMap(bencodedString string) (Node, error) {
	mapNode := MapNode{}
	remainingString := bencodedString[1:]

	for remainingString[0] != 'e' {
		key, err := decodeString(remainingString)
		if err != nil {
			return nil, err
		}
		remainingString = key.(ParseResult).RemainingString

		value, err := Decode(remainingString)
		if err != nil {
			return nil, err
		}
		remainingString = value.(ParseResult).RemainingString

		mapNode.Children[key.(ParseResult).Node.(StringNode).Data] = value
	}

	return mapNode, nil
}

func decodeString(bencodedString string) (Node, error) {
	var firstColonIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}

	return StringNode{bencodedString[firstColonIndex+1 : firstColonIndex+1+length]}, nil
}

func DecodeAll(bencodedString string) ([]Node, error) {
	var nodes []Node

	for len(bencodedString) > 0 {
		node, err := Decode(bencodedString)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
		bencodedString = node.(ParseResult).RemainingString
	}

	return nodes, nil
}

// now we need to implement the encode function

func Encode(node Node) (string, error) {
	switch node.(type) {
	case IntNode:
		return encodeInt(node.(IntNode))
	case StringNode:
		return encodeString(node.(StringNode))
	case ListNode:
		return encodeList(node.(ListNode))
	case MapNode:
		return encodeMap(node.(MapNode))
	default:
		return "", &InvalidBencodeError{"Invalid node"}
	}
}

func encodeInt(node IntNode) (string, error) {
	return "i" + strconv.Itoa(node.Data) + "e", nil
}

func encodeString(node StringNode) (string, error) {
	return strconv.Itoa(len(node.Data)) + ":" + node.Data, nil
}

func encodeList(node ListNode) (string, error) {
	encodedString := "l"

	for _, child := range node.Children {
		encodedChild, err := Encode(child)
		if err != nil {
			return "", err
		}
		encodedString += encodedChild
	}

	encodedString += "e"

	return encodedString, nil
}

func encodeMap(node MapNode) (string, error) {
	encodedString := "d"

	for key, child := range node.Children {
		encodedKey, err := encodeString(StringNode{key})
		if err != nil {
			return "", err
		}
		encodedChild, err := Encode(child)
		if err != nil {
			return "", err
		}
		encodedString += encodedKey + encodedChild
	}

	encodedString += "e"

	return encodedString, nil
}

func EncodeAll(nodes []Node) (string, error) {
	encodedString := ""

	for _, node := range nodes {
		encodedNode, err := Encode(node)
		if err != nil {
			return "", err
		}
		encodedString += encodedNode
	}

	return encodedString, nil
}