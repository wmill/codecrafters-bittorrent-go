package altbencode

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

type InvalidBencodeError struct {
	Message string
}

func (e *InvalidBencodeError) Error() string {
	return e.Message
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