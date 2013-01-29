package kilt

import (
	. "github.com/robertkrimen/terst"
	"testing"
)

func Test_StringTrimSpace(t *testing.T) {
	Terst(t)
	test := func(given, want string) {
		have := StringTrimSpace(given)
		Is(have, want)
	}
	test("  \n\n  Xyzzy", "  Xyzzy")
	test("  \n\n  Xyzzy\n", "  Xyzzy\n")
	test("  \n\n  ", "")
	test("\nXyzzy\nNothing happens.  \n\n", "Xyzzy\nNothing happens.\n")
	test("   ", "")
	test("Xyzzy\n  \n  \n", "Xyzzy\n")
	test("Xyzzy\n\n\n1", "Xyzzy\n\n\n1")
}
