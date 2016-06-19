package lib

import (
	"example.com/repo/lib/deep"
)

// Answer returns the ultimate answer to life, the universe and everything.
func Answer() int {
	var d deep.Thought
	return d.Compute()
}
