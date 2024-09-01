package data

import (
	"fmt"
	"testing"
)

func TestCreatePricingTale(t *testing.T) {
	s, err := NewPostgresStore()
	if err != nil {
		t.Fatal(err)
	}
	err = s.createPricingTable()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPricing(t *testing.T) {
	s, err := NewPostgresStore()
	id := 1122
	price, err := s.getPricing(id)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Pricing for id(%d):  %+v\n", id, price)

	ids, err := s.GetPlanIds()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Plan ids: %+v\n", ids)
}
