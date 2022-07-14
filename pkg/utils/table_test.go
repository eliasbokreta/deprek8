package utils

import (
	"testing"
)

func TestGetTableWriter(t *testing.T) {
	writer := GetTableWriter([]string{"a", "b", "c", "d"})

	writer.Append([]string{"1", "2", "3", "4"})
	writer.Append([]string{"1", "2", "3", "4"})
	writer.Append([]string{"1", "2", "3", "4"})

	if writer.NumLines() != 3 {
		t.Fatal("writer should have 3 lines")
	}

	writer.Render()
}
