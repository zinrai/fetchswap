package main

import "testing"

func TestModeIsParsedAsOctal(t *testing.T) {
	got, err := parseMode("0640")
	if err != nil || got != 0o640 {
		t.Errorf("parseMode(%q) = %o, %v, want 640", "0640", got, err)
	}
}

func TestSecondValueIsRefused(t *testing.T) {
	var s onceString
	if err := s.Set("a"); err != nil {
		t.Fatal(err)
	}
	if err := s.Set("b"); err == nil {
		t.Errorf("second Set succeeded, value %q", s.value)
	}
}
