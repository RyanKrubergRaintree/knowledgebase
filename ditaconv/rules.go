package ditaconv

import (
	"encoding/xml"
	"strings"

	"github.com/raintreeinc/knowledgebase/ditaconv/xmlconv"
)

// checks wheter dita tag corresponds to some "root element"
func isBodyTag(tag string) bool { return strings.Contains(tag, "body") }

// whether to process each child as separate item
func shouldUnwrap(name xml.Name) bool {
	switch name.Local {
	case "section",
		"example",
		"sectiondiv":
		return true
	}
	return false
}

func NewHTMLRules() *xmlconv.Rules {
	return &xmlconv.Rules{
		Translate: map[string]string{
			// conversion
			"xref": "a",
			"link": "a",

			//lists
			"choices":         "ul",
			"choice":          "li",
			"steps-unordered": "ul",
			"steps":           "ol",
			"step":            "li",
			"substeps":        "ol",
			"substep":         "li",

			"i":     "em",
			"lines": "pre",

			"codeblock": "code",

			"codeph":      "span",
			"cmdname":     "span",
			"cmd":         "span",
			"secright":    "span",
			"shortcut":    "span",
			"wintitle":    "span",
			"filepath":    "span",
			"menucascade": "span",

			"synph":    "span",
			"delim":    "span",
			"sep":      "span",
			"parmname": "span",

			"userinput": "kbd",

			"image": "img",

			// ui
			"uicontrol": "span",

			// divs
			"context":    "div",
			"result":     "div",
			"stepresult": "div",
			"stepxmp":    "div",
			"info":       "div",
			"note":       "div",
			"refsyn":     "div",
			"bodydiv":    "div",
			"fig":        "div", //TODO: convert to itemImage instead

			"prereq":  "div",
			"postreq": "div",

			// tables
			"simpletable": "table",
			"sthead":      "thead",
			"strow":       "tr",
			"stentry":     "td",

			"colspec": "colgroup",

			"row":   "tr",
			"entry": "td",
		},
		Remove: map[string]bool{
			"br":            true,
			"draft-comment": true,
			"colspec":       true,
		},
		Unwrap: map[string]bool{
			"dlentry": true,
			"tgroup":  true,
		},
		Callback: map[string]xmlconv.Callback{},
	}
}
