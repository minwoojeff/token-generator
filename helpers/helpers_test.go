package helpers

import "testing"

func TestOpen_Empty(t *testing.T) {
	_, err := Open("")
	if err == nil {
		t.Errorf("Expecting error for directory, got back nil")
	}
}

func TestOpen_NotExistent(t *testing.T) {
	_, err := Open("something")
	if err == nil {
		t.Error("Expecting error for non existent directory, got back nil")
	}
}
