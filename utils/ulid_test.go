package utils

import (
	"fmt"
	"testing"
)

func TestULID(t *testing.T) {
	id, err := GenerateULID()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("id : %s ", id)
	if len(id) > 26 || len(id) < 26 {
		t.Fatal("invalid ulid")
	}
}
