package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func compare(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Logf("\nexpected %q\n  actual %q", expected, actual)
		t.Fail()
	}
}

func test(t *testing.T, src string, expected string) {
	d := json.NewDecoder(strings.NewReader(src))
	ctx := newContext(d)
	err := ctx.findPaths("")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	out := ctx.fmtOutput(false, true)
	compare(t, out, expected)
}

func TestArray(t *testing.T)  { test(t, "[1, 2, 3]", "[] = 1\n") }
func TestObject(t *testing.T) { test(t, `{ "a": 1, "b": "two" }`, ".a = 1\n.b = \"two\"\n") }

func TestArrayInObject(t *testing.T) {
	test(t, `{
		"array": [1, 2, 3],
		"str": "hi"
	}`, ".array[] = 1\n.str = \"hi\"\n")
}

func TestObjectInArray(t *testing.T) {
	test(t, `[
		1,
		{ "a": "hi" }
	]`, "[] = 1\n[].a = \"hi\"\n")
}

func TestEverything(t *testing.T) {
	src := `{
		"id": 12,
		"name": "alligator",
		"active": true,
		"languages": [
			{
				"name": "JavaScript",
				"yoe": 8
			},
			{
				"name": "Go",
				"yoe": 2
			},
			{
				"name": "Python",
				"yoe": 13
			}
		]
	}`

	expected := ".id = 12\n" +
		".name = \"alligator\"\n" +
		".active = true\n" +
		".languages[].name = \"JavaScript\"\n" +
		".languages[].yoe = 8\n"

	test(t, src, expected)
}
