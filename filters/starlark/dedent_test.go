// https://github.com/lithammer/dedent/blob/master/dedent_test.go
package starlarkfilter

import (
	"fmt"
	"testing"
)

const errorMsg = "\nexpected %q\ngot %q"

type dedentTest struct {
	text, expect string
}

func TestDedentNoMargin(t *testing.T) {
	texts := []string{
		// No lines indented
		"Hello there.\nHow are you?\nOh good, I'm glad.",
		// Similar with a blank line
		"Hello there.\n\nBoo!",
		// Some lines indented, but overall margin is still zero
		"Hello there.\n  This is indented.",
		// Again, add a blank line.
		"Hello there.\n\n  Boo!\n",
	}

	for _, text := range texts {
		if text != dedent(text) {
			t.Errorf(errorMsg, text, dedent(text))
		}
	}
}

func TestDedentEven(t *testing.T) {
	texts := []dedentTest{
		{
			// All lines indented by two spaces
			text:   "  Hello there.\n  How are ya?\n  Oh good.",
			expect: "Hello there.\nHow are ya?\nOh good.",
		},
		{
			// Same, with blank lines
			text:   "  Hello there.\n\n  How are ya?\n  Oh good.\n",
			expect: "Hello there.\n\nHow are ya?\nOh good.\n",
		},
		{
			// Now indent one of the blank lines
			text:   "  Hello there.\n  \n  How are ya?\n  Oh good.\n",
			expect: "Hello there.\n\nHow are ya?\nOh good.\n",
		},
	}

	for _, text := range texts {
		if text.expect != dedent(text.text) {
			t.Errorf(errorMsg, text.expect, dedent(text.text))
		}
	}
}

func TestDedentUneven(t *testing.T) {
	texts := []dedentTest{
		{
			// Lines indented unevenly
			text: `
			def foo():
				while 1:
					return foo
			`,
			expect: `
def foo():
	while 1:
		return foo
`,
		},
		{
			// Uneven indentation with a blank line
			text:   "  Foo\n    Bar\n\n   Baz\n",
			expect: "Foo\n  Bar\n\n Baz\n",
		},
		{
			// Uneven indentation with a whitespace-only line
			text:   "  Foo\n    Bar\n \n   Baz\n",
			expect: "Foo\n  Bar\n\n Baz\n",
		},
	}

	for _, text := range texts {
		if text.expect != dedent(text.text) {
			t.Errorf(errorMsg, text.expect, dedent(text.text))
		}
	}
}

// dedent() should not mangle internal tabs.
func TestDedentPreserveInternalTabs(t *testing.T) {
	text := "  hello\tthere\n  how are\tyou?"
	expect := "hello\tthere\nhow are\tyou?"
	if expect != dedent(text) {
		t.Errorf(errorMsg, expect, dedent(text))
	}

	// Make sure that it preserves tabs when it's not making any changes at all
	if expect != dedent(expect) {
		t.Errorf(errorMsg, expect, dedent(expect))
	}
}

// dedent() should not mangle tabs in the margin (i.e. tabs and spaces both
// count as margin, but are *not* considered equivalent).
func TestDedentPreserveMarginTabs(t *testing.T) {
	texts := []string{
		"  hello there\n\thow are you?",
		// Same effect even if we have 8 spaces
		"        hello there\n\thow are you?",
	}

	for _, text := range texts {
		d := dedent(text)
		if text != d {
			t.Errorf(errorMsg, text, d)
		}
	}

	texts2 := []dedentTest{
		{
			// dedent() only removes whitespace that can be uniformly removed!
			text:   "\thello there\n\thow are you?",
			expect: "hello there\nhow are you?",
		},
		{
			text:   "  \thello there\n  \thow are you?",
			expect: "hello there\nhow are you?",
		},
		{
			text:   "  \t  hello there\n  \t  how are you?",
			expect: "hello there\nhow are you?",
		},
		{
			text:   "  \thello there\n  \t  how are you?",
			expect: "hello there\n  how are you?",
		},
	}

	for _, text := range texts2 {
		if text.expect != dedent(text.text) {
			t.Errorf(errorMsg, text.expect, dedent(text.text))
		}
	}
}

func Examplededent() {
	s := `
		Lorem ipsum dolor sit amet,
		consectetur adipiscing elit.
		Curabitur justo tellus, facilisis nec efficitur dictum,
		fermentum vitae ligula. Sed eu convallis sapien.`
	fmt.Println(dedent(s))
	fmt.Println("-------------")
	fmt.Println(s)
	// Output:
	// Lorem ipsum dolor sit amet,
	// consectetur adipiscing elit.
	// Curabitur justo tellus, facilisis nec efficitur dictum,
	// fermentum vitae ligula. Sed eu convallis sapien.
	// -------------
	//
	//		Lorem ipsum dolor sit amet,
	//		consectetur adipiscing elit.
	//		Curabitur justo tellus, facilisis nec efficitur dictum,
	//		fermentum vitae ligula. Sed eu convallis sapien.
}

func BenchmarkDedent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dedent(`Lorem ipsum dolor sit amet, consectetur adipiscing elit.
		Curabitur justo tellus, facilisis nec efficitur dictum,
		fermentum vitae ligula. Sed eu convallis sapien.`)
	}
}
