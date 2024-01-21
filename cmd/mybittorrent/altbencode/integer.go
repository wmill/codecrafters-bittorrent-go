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
	endIndex := startIndex;

	for bencodedString[endIndex] != 'e' {
		endIndex++
	}
	
	stringInt := bencodedString[startIndex+1:endIndex]
	integer, err := strconv.Atoi(stringInt)

	
	return ParseResult{
		Node: IntNode{
			Data: integer,
		},
		RemainingStringIndex: endIndex + 1,
	}, err
}

func encodeInteger(integer Node) (string, error) {
	return "i" + strconv.Itoa(integer.GetData().(int)) + "e", nil
}
