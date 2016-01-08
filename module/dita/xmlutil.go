package dita

import (
	"encoding/xml"
	"sort"
)

func getAttr(n *xml.StartElement, key string) (val string) {
	for _, attr := range n.Attr {
		if attr.Name.Local == key {
			return attr.Value
		}
	}
	return ""
}

type attrByName []xml.Attr

func (xs attrByName) Len() int           { return len(xs) }
func (xs attrByName) Swap(i, j int)      { xs[i], xs[j] = xs[j], xs[i] }
func (xs attrByName) Less(i, j int) bool { return xs[i].Name.Local < xs[j].Name.Local }

func setAttr(n *xml.StartElement, key, val string) {
	n.Attr = append([]xml.Attr{}, n.Attr...)

	for i := range n.Attr {
		attr := &n.Attr[i]
		if attr.Name.Local == key {
			if val == "" {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			} else {
				attr.Value = val
			}
			return
		}
	}

	if val == "" {
		return
	}

	n.Attr = append(n.Attr, xml.Attr{
		Name:  xml.Name{Local: key},
		Value: val,
	})
	sort.Sort(attrByName(n.Attr))
}
