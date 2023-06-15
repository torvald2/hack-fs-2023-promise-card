package polybase

import (
	"fmt"
	"testing"
)

func TestCreateRecord(t *testing.T) {
	key := ""
	cl, err := NewPolybaseClient("", "https://testnet.polybase.xyz")
	if err != nil {
		t.Fatal(err)
	}
	args := make([]interface{}, 0)
	args = append(args, "1")
	args = append(args, "2")
	data, err := cl.CreateRecord("User", args, key)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(data)

}

func TestGetRecord(t *testing.T) {
	key := ""
	cl, err := NewPolybaseClient("", "https://testnet.polybase.xyz")
	if err != nil {
		t.Fatal(err)
	}

	data, err := cl.GetRecord("User", key, "1")
	if err != nil {
		t.Error(err)
	}
	fmt.Print(data)

}
