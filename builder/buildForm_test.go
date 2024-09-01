package builder

import (
	"testing"

	"github.com/sonal3323/form-poc/data"
	"github.com/sonal3323/form-poc/types"
)

func TestRenderForm(t *testing.T) {
	fd := types.FormMetaData{
		Id:     "01HQFBPWRESFPPD28BEDEK8H2G",
		UserId: 5,
	}

	ps, err := data.NewPostgresStore()
	if err != nil {
		t.Fatal(err)
	}

	r := NewFormBuilder(*ps)
	if _, err := r.BuildForm(fd); err != nil {
		t.Fatal(err)
	}

}
