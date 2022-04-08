package main

import (
	"testing"
)

func TestCount(t *testing.T) {
	s := `
		if err == nil {}

		if err != nil {}

		int main() {
			return 0;
		}
	`

	if caculate(s) != 2 {
		t.Error("Expected 2, got", caculate(s))
	}
}
