package data

import (
	"fmt"
	"testing"
)

func TestQuestionTemplates(t *testing.T) {
	fmt.Println("-------- Question Template Test-------------")

	s, err := NewPostgresStore()
	if err := s.createQuestionTemplatesTable(); err != nil {
		t.Fatal(err)
	}

	templs, err := s.GetQuestionTemplates()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Templates :  %+v\n", templs)

	fmt.Printf("----------------end--------")
}
