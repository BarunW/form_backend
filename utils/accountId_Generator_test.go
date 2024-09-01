package utils

import (
	"fmt"
	"testing"
)

func TestGeneratorAccountId(t *testing.T) {
	id, err := GenerateAccountId()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(id)
	if len(id) != 11 {
		t.Fatal("lenth of id is should be 11")
	}
}
