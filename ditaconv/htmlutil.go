package ditaconv

import (
	"encoding/xml"

	"github.com/egonelbre/fedwiki"
)

func convertTags(keywords []string) []string {
	tags := []string{}
	for _, tag := range keywords {
		slug := string(fedwiki.Slugify(tag))
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
