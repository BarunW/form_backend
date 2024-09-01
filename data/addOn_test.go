package data

import (
	"fmt"
	"testing"

	"github.com/sonal3323/form-poc/types"
)

func TestGetIs(t *testing.T) {
	fmt.Println("=======================")
	s, err := NewPostgresStore()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("----------createResponseLimitTable()----------")
	err = s.createResponseLimitTable()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("----------GetAllResponseLimit()----------")
	data, err := s.GetAllResponseLimit()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Response Limit : %+v\n", data)
	fmt.Println("----------getResponseLimit(id)----------")
	var id types.ResponseLimitId = "250-BASIC-RESPLIMIT"
	d, err := s.getResponseLimit(id)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Response Limit : %+v\n", d)

	d2, err := s.getResponseLimit("")
	if err != nil && err.Error() != "404" {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Printf("Response Limit : %+v\n", d2)
}
