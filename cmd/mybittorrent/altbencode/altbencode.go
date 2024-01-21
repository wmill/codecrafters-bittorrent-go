package altbencode

import "strconv"

type Node interface {
	GetData() interface{}
}

type ParseResult struct {
	Node Node
	RemainingStringIndex int
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

func decode(bencodedString string, startIndex int) (ParseResult, error){
	if len(bencodedString) == 0 {
		return ParseResult{}, nil
	}

	switch bencodedString[startIndex] {
	case 'i':
		return decodeInteger(bencodedString, startIndex)
	case 'l':
		return decodeList(bencodedString, startIndex)
	case 'd':
		return decodeDictionary(bencodedString, startIndex)
	default:
		return decodeString(bencodedString, startIndex)
	}
}

func encode(node Node) (string, error) {
	switch node.(type) {
	case StringNode:
		return encodedString(node)
	case IntNode:
		return encodeInteger(node)
	case ListNode:
		return encodeList(node)
	case MapNode:
		return encodeDictionary(node)
	default:
		return "", &InvalidBencodeError{"Unknown node type"}
	}
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

func decodeInteger(bencodedString string, startIndex int) (ParseResult, error) {
	length := len(bencodedString)
	if bencodedString[length-1] != 'e' {
		return ParseResult{}, &InvalidBencodeError{"Invalid integer"}
	}
	stringInt := bencodedString[startIndex+1:length-1]
	integer, err := strconv.Atoi(stringInt)
	return ParseResult{
		Node: IntNode{
			Data: integer,
		},
		RemainingStringIndex: length,
	}, err
}

func encodeInteger(integer Node) (string, error) {
	return "i" + strconv.Itoa(integer.GetData().(int)) + "e", nil
}

func decodeList(bencodedString string, startIndex int) (ParseResult, error) {
	var children []Node
	index := startIndex + 1
	for index < len(bencodedString) {
		if bencodedString[index] == 'e' {
			break
		}
		result, err := decode(bencodedString, index)
		if err != nil {
			return ParseResult{}, err
		}
		children = append(children, result.Node)
		index = result.RemainingStringIndex
	}
	return ParseResult{
		Node: ListNode{
			Children: children,
		},
		RemainingStringIndex: index + 1,
	}, nil
}


func encodeList(list Node) (string, error) {
	var result string
	result += "l"
	for _, node := range list.GetData().([]Node) {
		encodedNode, err := encode(node)
		if err != nil {
			return "", err
		}
		result += encodedNode
	}
	result += "e"
	return result, nil
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
	for key, node := range dictionary.GetData().(map[string]Node) {
		encodedKey, err := encodedString(StringNode{key})
		if err != nil {
			return "", err
		}
		encodedNode, err := encode(node)
		if err != nil {
			return "", err
		}
		result += encodedKey + encodedNode
	}
	result += "e"
	return result, nil
}

func Decode(bencodedString string) (Node, error) {
	result, err := decode(bencodedString, 0)
	if err != nil {
		return nil, err
	}
	return result.Node, nil
}

func Encode(node Node) (string, error) {
	return encode(node)
}