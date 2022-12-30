package interpreter

import (
	"testing"
)

func TestStringConcat(t *testing.T) {
	testWithOutput(
		`a := "hello world"; b := a.concat("!"); noot!(a); noot!(b)`,
		"hello world\nhello world!\n",
		t,
	)
}
