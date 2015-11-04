// +build functional

package main_test

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

const binary = "bin.test"

var tests = []struct {
	name string
	args []string
	yes  []string // Regular expressions that should match.
	no   []string // Regular expressions that should not match.
}{
	// Sanity check.
	{
		"sanity check",
		[]string{"world"},
		[]string{`Hello world`},
		nil,
	},
}

func TestMain(m *testing.M) {
	flag.Parse()

	// go build
	out, err := exec.Command("go", "build", "-o", binary).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "building failed: %v\n%s", err, out)
		os.Exit(2)
	}

	r := m.Run()
	os.Remove(binary)
	os.Exit(r)
}

func TestFoo(t *testing.T) {
	for _, test := range tests {
		cmd := exec.Command("./"+binary, test.args...)
		output := run(cmd, t)
		failed := false
		for j, yes := range test.yes {
			re, err := regexp.Compile(yes)
			if err != nil {
				t.Fatalf("%s.%d: compiling %#q: %s", test.name, j, yes, err)
			}
			if !re.Match(output) {
				t.Errorf("%s.%d: no match for %s %#q", test.name, j, test.args, yes)
				failed = true
			}
		}
		for j, no := range test.no {
			re, err := regexp.Compile(no)
			if err != nil {
				t.Fatalf("%s.%d: compiling %#q: %s", test.name, j, no, err)
			}
			if re.Match(output) {
				t.Errorf("%s.%d: incorrect match for %s %#q", test.name, j, test.args, no)
				failed = true
			}
		}
		if failed {
			t.Logf("\n%s", output)
		}
	}
}

// run runs the command, but calls t.Fatal if there is an error.
func run(c *exec.Cmd, t *testing.T) []byte {
	output, err := c.CombinedOutput()
	if err != nil {
		os.Stdout.Write(output)
		t.Fatal(err)
	}
	return output
}
