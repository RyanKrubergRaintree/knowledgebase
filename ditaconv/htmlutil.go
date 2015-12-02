package ditaconv

import (
	"encoding/xml"

	"github.com/raintreeinc/knowledgebase/ditaconv/dita"
	"github.com/raintreeinc/knowledgebase/kb"
)

func convertTags(prolog *dita.Prolog) []string {
	tags := []string{}

	raw := []string{}
	raw = append(raw, prolog.Keywords...)
	for _, rid := range prolog.ResourceID {
		raw = append(raw, "id/"+rid.Name)
	}

	for _, tag := range raw {
		slug := string(kb.Slugify(tag))
		if slug == "" || slug == "-" {
			continue
		}

		contains := false
		for _, tag := range tags {
			if tag == string(slug) {
				contains = true
				break
			}
		}
		if contains {
			continue
		}

		tags = append(tags, string(slug))
	}

	return tags
}

func asParagraph(v string) (string, bool) {
	// TODO
	return v, false
}

func getAttr(n *xml.StartElement, key string) (val string) {
	for _, attr := range n.Attr {
		if attr.Name.Local == key {
			return attr.Value
		}
	}
	return ""
}

func setAttr(n *xml.StartElement, key, val string) {
	for i, attr := range n.Attr {
		if attr.Name.Local == key {
			n.Attr[i].Value = val
			return
		}
	}
	n.Attr = append(n.Attr, xml.Attr{
		Name:  xml.Name{Local: key},
		Value: val,
	})
}
