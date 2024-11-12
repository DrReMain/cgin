package json

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Likes []string `json:"likes"`
}

func TestMarshalUnmarshal(t *testing.T) {
	original := testStruct{
		Name:  "Alice",
		Age:   30,
		Likes: []string{"Coding", "Reading"},
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded testStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Decoded value doesn't match original. Got %+v, want %+v", decoded, original)
	}
}

func TestMarshalToString(t *testing.T) {
	original := testStruct{
		Name:  "Bob",
		Age:   25,
		Likes: []string{"Music", "Sports"},
	}

	str, err := MarshalToString(original)
	if err != nil {
		t.Fatalf("MarshalToString failed: %v", err)
	}

	expected := `{"name":"Bob","age":25,"likes":["Music","Sports"]}`
	if str != expected {
		t.Errorf("MarshalToString produced unexpected result.\nGot:  %s\nWant: %s", str, expected)
	}

	var decoded testStruct
	err = Unmarshal([]byte(str), &decoded)
	if err != nil {
		t.Fatalf("Unmarshal of MarshalToString result failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Decoded value from MarshalToString doesn't match original. Got %+v, want %+v", decoded, original)
	}
}

func TestMarshalEmptyStruct(t *testing.T) {
	empty := testStruct{}

	data, err := Marshal(empty)
	if err != nil {
		t.Fatalf("Marshal empty struct failed: %v", err)
	}

	expected := `{"name":"","age":0,"likes":null}`
	if string(data) != expected {
		t.Errorf("Marshal empty struct produced unexpected result.\nGot:  %s\nWant: %s", string(data), expected)
	}

	var decoded testStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal of empty struct failed: %v", err)
	}

	if !reflect.DeepEqual(empty, decoded) {
		t.Errorf("Decoded empty struct doesn't match original. Got %+v, want %+v", decoded, empty)
	}
}
