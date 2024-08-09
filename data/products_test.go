package data

import "testing"

func TestChecksValidation(t *testing.T) {
	prod := &Product{
		Name:  "test",
		Price: 1,
		SKU:   "abc-abc-abc",
	}

	err := prod.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
