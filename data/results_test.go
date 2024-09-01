package data_test

import (
	"testing"

	"github.com/sonal3323/form-poc/data"
)

func TestGetAnswer(t *testing.T) {
	np, err := data.NewPostgresStore()
	if err != nil {
		t.Fatal(err)
	}

	formId := "01HSGTF548E0FJFT4VGX6YZB5E"

	_, err = np.GetResult(formId)
	if err != nil {
		t.Fatal(err)
	}

}
