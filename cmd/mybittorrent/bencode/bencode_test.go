package bencode

import (
	"reflect"
	"testing"
)

func TestDecodeAll(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    []Node
        wantErr bool
    }{
        {
            name:    "empty string",
            input:   "",
            want:    []Node{},
            wantErr: false,
        },
        // Add more test cases here
        // {
        // 	name:    "test case 2",
        // 	input:   "your input here",
        // 	want:    []Node{your expected nodes here},
        // 	wantErr: false or true based on your expectation,
        // },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := DecodeAll(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("DecodeAll() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("DecodeAll() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestDeocdeAndEncode(t *testing.T) {
	testString := "l4:spam4:eggse"
	// testString := "d3:cow3:moo4:spam4:eggse"

	nodes, err := Decode(testString)
	if err != nil {
		t.Error(err)
	}
	resultString, err2 := Encode(nodes)

	if err2 != nil {
		t.Error(err2)
	}

	if resultString != testString {
		t.Errorf("Expected %s, got %s", testString, resultString)
	}
}