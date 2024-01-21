package altbencode

import "testing"

func TestDecodeAndEncodeString(t *testing.T) {
	testString := "4:spam"
	result, err := decodeString(testString,0)
	if err != nil {
		t.Error(err)
	}
	nodes := result.Node
	remainingStringIndex := result.RemainingStringIndex

	if remainingStringIndex != len(testString) {
		t.Error("Expected remainingStringIndex to be equal to length of testString")
	}

	resultString, err2 := encodedString(nodes)

	if err2 != nil {
		t.Error(err2)
	}

	if resultString != testString {
		t.Errorf("Expected %s, got %s", testString, resultString)
	}
}

func TestDecodeAndEncodeInteger(t *testing.T) {
	testString := "i3e"
	result, err := decodeInteger(testString,0)
	if err != nil {
		t.Error(err)
	}

	nodes := result.Node
	remainingStringIndex := result.RemainingStringIndex

	if remainingStringIndex != len(testString) {
		t.Error("Expected remainingStringIndex to be equal to length of testString")
	
	}
	resultString, err2 := encodeInteger(nodes)

	if err2 != nil {
		t.Error(err2)
	}

	if resultString != testString {
		t.Errorf("Expected %s, got %s", testString, resultString)
	}
}


func TestDecodeAndEncodeList(t *testing.T) {
	testString := "l4:spam4:eggse"
	result, err := decodeList(testString,0)
	if err != nil {
		t.Error(err)
	}
	nodes := result.Node
	remainingStringIndex := result.RemainingStringIndex

	if remainingStringIndex != len(testString) {
		t.Error("Expected remainingStringIndex to be equal to length of testString")
	
	}

	resultString, err2 := encodeList(nodes)

	if err2 != nil {
		t.Error(err2)
	}

	if resultString != testString {
		t.Errorf("Expected %s, got %s", testString, resultString)
	}
}

func TestDecodeAndEncodeDictionary(t *testing.T) {
	testString := "d3:cow3:moo4:spam4:eggse"
	result, err := decodeDictionary(testString,0)
	if err != nil {
		t.Error(err)
	}
	nodes := result.Node
	remainingStringIndex := result.RemainingStringIndex

	if remainingStringIndex != len(testString) {
		t.Error("Expected remainingStringIndex to be equal to length of testString")
	
	}

	resultString, err2 := encodeDictionary(nodes)

	if err2 != nil {
		t.Error(err2)
	}

	if resultString != testString {
		t.Errorf("Expected %s, got %s", testString, resultString)
	}
}