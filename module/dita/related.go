package dita

import (
	"github.com/raintreeinc/ditaconvert"
	"github.com/raintreeinc/ditaconvert/dita"
	"github.com/raintreeinc/ditaconvert/html"
)

func (conversion *PageConversion) RelatedLinksAsHTML() (div string) {
	context := conversion.Context

	topic := context.Topic
	if topic == nil || ditaconvert.EmptyLinkSets(topic.Links) {
		return ""
	}

	contains := func(xs []string, s string) bool {
		for _, x := range xs {
			if x == s {
				return true
			}
		}
		return false
	}

	div += `<div class="related-links">`

	for _, set := range topic.Links {
		if len(set.Children) > 0 {
			if set.CollType == dita.Sequence {
				div += "<ol>"
			} else {
				div += "<ul>"
			}

			for _, link := range set.Children {
				div += "<li>" + conversion.LinkAsAnchor(link)
				if link.Topic.Synopsis != "" {
					div += "<p>" + link.Topic.Synopsis + "</p>"
				}
				div += "</li>"
			}

			if set.CollType == dita.Sequence {
				div += "</ol>"
			} else {
				div += "</ul>"
			}
		}

		if set.Parent != nil {
			div += "<div><b>Parent topic: </b>" + conversion.LinkAsAnchor(set.Parent) + "</div>"
		}
		if set.Prev != nil {
			div += "<div><b>Previous topic: </b>" + conversion.LinkAsAnchor(set.Prev) + "</div>"
		}
		if set.Next != nil {
			div += "<div><b>Next topic: </b>" + conversion.LinkAsAnchor(set.Next) + "</div>"
		}
	}

	grouped := make(map[string][]*ditaconvert.Link)
	order := []string{"tutorial", "concept", "task", "reference", "information"}
	for _, set := range topic.Links {
		for _, link := range set.Siblings {
			kind := ""
			if link.Topic != nil {
				kind = link.Topic.Original.XMLName.Local
			}
			if link.Type != "" {
				kind = link.Type
			}
			if !contains(order, kind) {
				kind = "information"
			}

			grouped[kind] = append(grouped[kind], link)
		}
	}

	for _, kind := range order {
		links := grouped[kind]
		if len(links) == 0 {
			continue
		}

		if kind != "information" && len(links) > 1 {
			kind += "s"
		}
		div += "<div><b>Related " + kind + "</b>"
		for _, link := range links {
			div += "<div>" + conversion.LinkAsAnchor(link) + "</div>"
		}
		div += "</div>"
	}

	div += "</div>"

	return div
}

func (conversion *PageConversion) LinkAsAnchor(link *ditaconvert.Link) string {
	title := html.EscapeCharData(link.FinalTitle())
	if link.Scope == "external" {
		return `<a href="` + html.NormalizeURL(link.Href) + `" class="external-link" target="_blank" rel="nofollow">` + title + `</a>`
	}

	if link.Topic == nil {
		return `<span style="background: #f00">` + title + `</span>`
	}

	slug, ok := conversion.Mapping.ByTopic[link.Topic]
	if !ok {
		return `<span style="background: #f00">` + title + `</span>`
	}

	return `<a href="` + string(slug) + `" data-link="` + string(slug) + `">` + title + `</a>`
}
