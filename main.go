package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/elves/upgrade-scripts-for-0.11/parse"
)

func main() {
	filenames := os.Args[1:]
	if len(filenames) == 0 {
		fixed, ok := fix("[stdin]", func() ([]byte, error) {
			return ioutil.ReadAll(os.Stdin)
		})
		if ok {
			os.Stdout.Write(fixed)
		}
	} else {
		// Fix files.
		for _, filename := range filenames {
			fixed, ok := fix(filename, func() ([]byte, error) {
				return ioutil.ReadFile(filename)
			})
			if ok {
				ioutil.WriteFile(filename, fixed, 0644)
			}
		}
	}
}

func fix(filename string, readall func() ([]byte, error)) ([]byte, bool) {
	src, err := readall()
	if err != nil {
		log.Printf("cannot read %s: %s; skipping", filename, err)
		return nil, false
	}
	if !utf8.Valid(src) {
		log.Printf("%s not utf8; skipping", filename)
		return nil, false
	}
	// os.Stdout.WriteString(fix("[stdin]", src))
	chunk, err := parse.Parse(filename, string(src))
	if err != nil {
		log.Printf("cannot parse %s: %s; skipping", filename, err)
		return nil, false
	}

	buf := new(bytes.Buffer)
	fixNode(chunk, buf)
	return buf.Bytes(), true
}

var (
	suppress = map[string]bool{
		"\n": true, ";": true, "in": true,
		"do": true, "done": true,
		"then": true, "fi": true,
		"tried": true,
	}
	addSpace = map[string]bool{
		"elif": true, "else": true,
		"except": true, "finally": true,
	}
)

func fixNode(n parse.Node, w io.Writer) {
	switch n := n.(type) {
	case *parse.Primary:
		if n.Type == parse.Variable {
			// $&x
			// to
			// $x~
			w.Write([]byte("$" + convert(n.Value)))
		} else {
			fixNodeDefault(n, w)
		}
	case *parse.Assignment:
		// '&x'={ } echo
		// to
		// x~={ } echo
		for _, child := range n.Children() {
			if child == n.Left {
				fixLHS(n.Left, w)
			} else {
				fixNode(child, w)
			}
		}
	case *parse.Form:
		for _, child := range n.Children() {
			if parse.IsCompound(child) {
				var varChild *parse.Compound
				for _, v := range n.Vars {
					if v == child {
						varChild = v
						break
					}
				}
				if varChild != nil {
					for _, grandchild := range varChild.Children() {
						if in := parse.GetIndexing(grandchild); in != nil {
							fixLHS(in, w)
							continue
						}
						fixNode(grandchild, w)
					}
					continue
				}
			}
			fixNode(child, w)
		}
	case *parse.Compound:
		// $x~foo
		// to
		// $x''~foo
		firstIndexing := true
		for _, child := range n.Children() {
			if parse.IsIndexing(child) {
				if !firstIndexing && strings.HasPrefix(child.SourceText(), "~") {
					w.Write([]byte("''"))
				}
				firstIndexing = false
			}
			fixNode(child, w)
		}
	default:
		fixNodeDefault(n, w)
	}
}

func fixNodeDefault(n parse.Node, w io.Writer) {
	if len(n.Children()) == 0 {
		text := n.SourceText()
		w.Write([]byte(text))
	} else {
		for _, child := range n.Children() {
			fixNode(child, w)
		}
	}
}

func fixLHS(n *parse.Indexing, w io.Writer) {
	for _, child := range n.Children() {
		if child == n.Head {
			w.Write([]byte(convert(n.Head.Value)))
		} else {
			fixNode(child, w)
		}
	}
}

func convert(old string) string {
	ns, name := parseVariableQName(old)
	if strings.HasPrefix(name, "&") {
		name = name[1:] + "~"
	}
	qname := ns + name
	if strings.ContainsRune(qname, '&') {
		fmt.Fprintf(os.Stderr, "Warning: rewritten variable $%s still contains &\n", qname)
	}
	return qname
}

func parseVariableQName(qname string) (ns, name string) {
	i := strings.LastIndexByte(qname, ':')
	if i == -1 {
		return "", qname
	}
	return qname[:i+1], qname[i+1:]
}
