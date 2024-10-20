package golang

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/sqlc-gen-go/internal/opts"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Struct struct {
	Table   *plugin.Identifier
	Name    string
	Fields  []Field
	Comment string
}

// Initialize title caser
var caser = cases.Title(language.English)

func StructName(name string, options *opts.Options) string {
	if rename := options.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		if unicode.IsDigit(r) {
			return r
		}
		return rune('_')
	}, name)

	if options.KeepCase {
		// If KeepCase is set, return the name as is.
		out = name
	} else {
		// Split the name by underscores and capitalize each part.
		for _, p := range strings.Split(name, "_") {
			if _, found := options.InitialismsMap[p]; found {
				out += strings.ToUpper(p)
			} else {
				out += caser.String(p)
			}
		}
	}

	// If a name has a digit as its first char, prepend an underscore to make it a valid Go name.
	r, _ := utf8.DecodeRuneInString(out)
	if unicode.IsDigit(r) {
		return "_" + out
	} else {
		return out
	}
}
