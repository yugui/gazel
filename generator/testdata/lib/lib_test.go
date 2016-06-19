package lib

import (
	"testing"
)

func TestAnswer(t *testing.T) {
	if got, want := Answer(), 42; got != want {
		t.Errorf("Answer() = %d; want %d", got, want)
	}
}
