package altbencode

import "encoding/json"

type ListNode struct {
	Children []Node
}

func (n ListNode) GetData() interface{} {
	return n.Children
}

//add a MarshalJSON method to ListNode
//add a MarshalJSON method 
func (n ListNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Children)
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
