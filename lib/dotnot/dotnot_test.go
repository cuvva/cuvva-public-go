package dotnot

import (
	"encoding/json"
	"testing"
)

var testsuite = []struct {
	name string
	from string
	to   string
}{
	{
		name: "Do not change single root key",
		from: `{"root":"key"}`,
		to:   `{"root":"key"}`,
	},
	{
		name: "Do not change multiple root keys",
		from: `{"another":"key","root":"key"}`,
		to:   `{"another":"key","root":"key"}`,
	},
	{
		name: "Dotnotate a single level",
		from: `{"nesting":{"one":true},"root":"key"}`,
		to:   `{"nesting.one":true,"root":"key"}`,
	},
	{
		name: "Dotnotate a bunch of levels",
		from: `{"nest":{"two":{"three":{"four":"layers_deep"}},"two_half":{"and_siblings_too":true}},"root":"key"}`,
		to:   `{"nest.two.three.four":"layers_deep","nest.two_half.and_siblings_too":true,"root":"key"}`,
	},
	{
		name: "Dotnotate arrays safely",
		from: `{"nest":{"two":["arrays","must","make","it","through"],"two_half":{"and_siblings_too":true}},"root":"key"}`,
		to:   `{"nest.two":["arrays","must","make","it","through"],"nest.two_half.and_siblings_too":true,"root":"key"}`,
	},
	{
		name: "Dotnotate array of objects",
		from: `{"nest":{"two":[{"key":"one"},{"key":"two"}]},"root":"key"}`,
		to:   `{"nest.two.0.key":"one","nest.two.1.key":"two","root":"key"}`,
	},
	{
		name: "Dotnotate set of arrays",
		from: `{"one":[{"set":[1,2,3]},{"set":[1,2,3]},{"set":[1,2,3]}]}`,
		to:   `{"one.0.set":[1,2,3],"one.1.set":[1,2,3],"one.2.set":[1,2,3]}`,
	},
}

func TestTo(t *testing.T) {
	for _, test := range testsuite {
		var v map[string]interface{}
		if err := json.Unmarshal([]byte(test.from), &v); err != nil {
			t.Error(err)
		}

		o := To(v)

		out, err := json.Marshal(o)
		if err != nil {
			t.Error(err)
		}

		if string(out) != test.to {
			t.Errorf("Test \"%s\" failed\nGot: %s\nExpected: %s\n", test.name, string(out), test.to)
		}
	}
}

func TestFrom(t *testing.T) {
	for _, test := range testsuite {
		var v map[string]interface{}
		if err := json.Unmarshal([]byte(test.to), &v); err != nil {
			t.Error(err)
		}

		o := From(v)

		out, err := json.Marshal(o)
		if err != nil {
			t.Error(err)
		}

		if string(out) != test.from {
			t.Errorf("Test \"%s\" failed\nGot: %s\nExpected: %s\n", test.name, string(out), test.from)
		}
	}
}
