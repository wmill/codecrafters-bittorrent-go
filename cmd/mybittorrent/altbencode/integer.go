package altbencode

import "strconv"

type IntNode struct {
	Data int
}

func (n IntNode) GetData() interface{} {
	return n.Data
}

//add a MarshalJSON method to IntNode
func (n IntNode) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(n.GetData().(int))), nil
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
