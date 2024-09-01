package data

import (
	"log"
	"testing"
)

func TestCreateTable(t *testing.T) {
	s, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	err = s.createSubscriptionTable()
	if err != nil {
		log.Fatal(err)
	}
}
