package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type pathEntry struct {
	path         string
	exampleValue interface{}
}

type context struct {
	paths []pathEntry
	seen  map[string]bool
	dec   *json.Decoder
	tok   *json.Token
	more  bool
	err   error
}

func newContext(d *json.Decoder) context {
	ctx := context{
		make([]pathEntry, 0),
		make(map[string]bool),
		d,
		nil,
		false,
		nil,
	}
	ctx.next()
	return ctx
}

func (ctx *context) next() *json.Token {
	if ctx.err != nil {
		return nil
	}

	ctx.more = ctx.dec.More()
	tok, err := ctx.dec.Token()
	if err == io.EOF {
		ctx.tok = nil
		return nil
	}
	if err != nil {
		ctx.err = err
		return nil
	}
	prev := ctx.tok
	ctx.tok = &tok
	return prev
}

func (ctx *context) add(path string, exampleValue interface{}) {
	if ctx.err != nil {
		return
	}

	if _, present := ctx.seen[path]; !present {
		ctx.paths = append(ctx.paths, pathEntry{path, exampleValue})
		ctx.seen[path] = true
	}
}

func (ctx *context) fmtOutput(colour bool, printValues bool) string {
	var sb strings.Builder
	for _, path := range ctx.paths {
		if !printValues {
			sb.WriteString(path.path + "\n")
			continue
		}

		v, err := json.Marshal(path.exampleValue)
		if err != nil {
			panic(err)
		}
		s := string(v)
		if len(s) > 64 {
			s = s[:64] + "..."
		}

		if colour {
			sb.WriteString(fmt.Sprintf("%s = \x1b[93m%s\x1b[0m\n", path.path, s))
		} else {
			sb.WriteString(fmt.Sprintf("%s = %s\n", path.path, s))
		}
	}
	return sb.String()
}

func (ctx *context) findPaths(path string) error {
	switch *ctx.tok {
	case json.Delim('['):
		ctx.next() // eat the [
		if ctx.err != nil {
			return ctx.err
		}

		for ctx.more {
			switch *ctx.tok {
			case json.Delim('{'), json.Delim('['):
				ctx.findPaths(path + "[]")
			default:
				ctx.add(path+"[]", *ctx.tok)
				ctx.next()
				if ctx.err != nil {
					return ctx.err
				}
			}
		}
		ctx.next() // eat the ]
	case json.Delim('{'):
		ctx.next() // eat the {
		for ctx.more {
			s := ctx.next()
			if ctx.err != nil {
				return ctx.err
			}

			if str, ok := (*s).(string); ok {
				switch *ctx.tok {
				case json.Delim('{'), json.Delim('['):
					ctx.findPaths(path + "." + str)
				default:
					ctx.add(path+"."+str, *ctx.tok)
					ctx.next()
					if ctx.err != nil {
						return ctx.err
					}
				}
			} else {
				return fmt.Errorf("expected a string but got %q", *s)
			}
		}
		ctx.next() // eat the }
	}

	return ctx.err
}

func main() {
	colour := flag.Bool("c", false, "use colour")
	printValues := flag.Bool("v", false, "show example values")
	flag.Parse()

	var d *json.Decoder
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		d = json.NewDecoder(f)
	} else {
		d = json.NewDecoder(os.Stdin)
	}

	ctx := newContext(d)
	err := ctx.findPaths("")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(ctx.fmtOutput(*colour, *printValues))
}
