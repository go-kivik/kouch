package testy

import "testing"

func TestTableTests(t *testing.T) {
	addFunc := func(a, b int) int {
		return a + b
	}
	type ttTest struct {
		a, b   int
		output int
	}
	table := NewTable()
	table.Add("one", func(_ *testing.T) interface{} {
		return ttTest{a: 1, output: 1}
	})
	table.Add("two", ttTest{a: 1, b: 1, output: 2})
	table.Add("three", func() interface{} {
		return ttTest{a: 1, b: 2, output: 3}
	})
	table.Run(t, func(t *testing.T, test ttTest) {
		output := addFunc(test.a, test.b)
		if output != test.output {
			t.Errorf("Expected %d, got %d\n", test.output, output)
		}
	})
}
