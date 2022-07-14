package utils

import "testing"

type OutputTest struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func TestOutput(t *testing.T) {
	outputTest := OutputTest{
		Field1: "a",
		Field2: "b",
	}

	err := OutputResult(outputTest, "json")
	if err != nil {
		t.Fatalf("Should not fail: %v", err)
	}

	err = OutputResult(outputTest, "yaml")
	if err != nil {
		t.Fatalf("Should not fail: %v", err)
	}

	err = OutputResult(outputTest, "")
	if err == nil {
		t.Fatal("Should fail")
	}

	err = OutputResult(make(chan int), "json")
	if err == nil {
		t.Fatal("Should fail")
	}
}
