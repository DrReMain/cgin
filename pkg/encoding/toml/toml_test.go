package toml

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Name  string
	Age   int
	Likes []string
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

	expected := `Name = "Bob"
Age = 25
Likes = ["Music", "Sports"]
`
	if str != expected {
		t.Errorf("MarshalToString produced unexpected result. Got:\n%s\nWant:\n%s", str, expected)
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
