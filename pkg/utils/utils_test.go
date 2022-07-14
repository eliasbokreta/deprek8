package utils

import "testing"

func TestStringInSlice(t *testing.T) {
	baseList := []string{"a", "b", "c", "d"}
	if StringInSlice("a", baseList) == false {
		t.Fatalf("Should return true")
	}
	if StringInSlice("ez", baseList) == true {
		t.Fatalf("Should return false")
	}
}
