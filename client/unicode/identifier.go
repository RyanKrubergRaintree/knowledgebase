// This program creates the support.js file
// with
//   go run support.go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"unicode"
)

type Range struct{ From, To rune }

func main() {
	ranges := []Range{}
	for r := rune(0); r < 0xFFFFFF; r += 1 {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			start := r
			for unicode.IsLetter(r) || unicode.IsNumber(r) {
				r += 1
			}
			ranges = append(ranges, Range{start, r})
		}
	}

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "// DO NOT MODIFY\n")
	fmt.Fprintf(&buf, "// THIS IS A GENERATED FILE\n\n")
	fmt.Fprintf(&buf, "package(\"kb.unicode\", function(exports){\n")
	fmt.Fprintf(&buf, "\t\"use strict\";\n\n")
	fmt.Fprintf(&buf, "\texports.IsIdent = IsIdent;\n")
	fmt.Fprintf(&buf, "\tfunction IsIdent(v){\n")
	fmt.Fprintf(&buf, "\t\tvar r = v.charCodeAt(0);\n")
	fmt.Fprintf(&buf, "\t\tfor(var i = 0; i < ident.length; i += 2){\n")
	fmt.Fprintf(&buf, "\t\t\tif((ident[i] <= r) && (r < ident[i+1])){\n")
	fmt.Fprintf(&buf, "\t\t\t\treturn true;\n")
	fmt.Fprintf(&buf, "\t\t\t}\n")
	fmt.Fprintf(&buf, "\t\t}\n")
	fmt.Fprintf(&buf, "\t\treturn false;\n")
	fmt.Fprintf(&buf, "\t}\n\n\n")

	fmt.Fprintf(&buf, "\tvar ident = [")
	for i, r := range ranges {
		if i > 0 {
			fmt.Fprintf(&buf, ", ")
		}
		if i%4 == 0 {
			fmt.Fprintf(&buf, "\n")
			fmt.Fprintf(&buf, "\t\t")
		}
		fmt.Fprintf(&buf, "0x%06x, 0x%06x", r.From, r.To)
	}
	fmt.Fprintf(&buf, "\n\t];\n")
	fmt.Fprintf(&buf, "});\n")

	ioutil.WriteFile("identifier.js", buf.Bytes(), 0755)
}
