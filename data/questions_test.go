package data

import (
	"fmt"
	"testing"
)

func TestGetQuestionTemplate(t *testing.T) {
	fmt.Println("=======================")
	s, err := NewPostgresStore()
	if err != nil {
		t.Fatal(err)
	}

	questTemp, err := s.GetQuestionTemplate("01HNFBT99TK7PA9B8KSHBS3NH8")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v", questTemp)

	questTemp2, err := s.GetQuestionTemplate("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", questTemp2)

}
