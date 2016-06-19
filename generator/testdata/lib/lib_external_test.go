package lib_test

import (
	"testing"

	"example.com/repo/lib"
)

func TestAnswerExternal(t *testing.T) {
	if got, want := lib.Answer(), 42; got != want {
		t.Errorf("lib.Answer() = %d; want %d", got, want)
	}
}
