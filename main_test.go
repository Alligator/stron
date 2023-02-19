package main

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func TestFiles(t *testing.T) {
	files, err := filepath.Glob("tests/*.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Log(file)

		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		expected, err := os.ReadFile(strings.Replace(file, ".json", ".stron", 1))
		if err != nil {
			t.Fatal(err)
		}

		test(t, string(src), string(expected))
	}
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

func TestNullExampleValues(t *testing.T) {
	src := `[
		{
			"alwaysNull": null,
			"sometimesNull": null
		},
		{
			"alwaysNull": null,
			"sometimesNull": "hello"
		},
		{
			"alwaysNull": null,
			"sometimesNull": "goodbye"
		}
	]`

	expected := "[].alwaysNull = null\n" +
		"[].sometimesNull = \"hello\"\n"

	test(t, src, expected)
}
