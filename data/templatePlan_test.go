package data

import (
	"fmt"
	"testing"
)

func TestCreateTemplPlanTable(t *testing.T) {
	fmt.Println("=============TemplatePlan===============")
	ps, err := NewPostgresStore()
	if err != nil {
		t.Fatal(err)
	}
	err = ps.createTemplPlanTable()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("--------------GET templPlan --------------")
	templPlanIds, err := ps.GetTemplatePlan(1121)
	if err != nil {
		t.Fatal(err)
	}

	if len(*templPlanIds) < 1 {
		t.Fatal(err)
	}
}
