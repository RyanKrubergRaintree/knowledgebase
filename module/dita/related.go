package dita

import (
	"github.com/raintreeinc/ditaconvert"
	"github.com/raintreeinc/ditaconvert/dita"
	"github.com/raintreeinc/ditaconvert/html"
)

func (conversion *PageConversion) RelatedLinksAsHTML() (div string) {
	topic := conversion.Topic
	if topic == nil || ditaconvert.EmptyLinkSets(topic.Links) {
		return "<div></div>"
	}

	contains := func(xs []string, s string) bool {
		for _, x := range xs {
			if x == s {
				return true
			}
		}
		return false
	}

	div += `<div>`

	var hasFamilyLinks bool
	for _, set := range topic.Links {
		hasFamilyLinks = hasFamilyLinks || set.Parent != nil || set.Prev != nil || set.Next != nil
		if len(set.Children) == 0 {
			continue
		}

		if set.CollType == dita.Sequence {
			div += `<ol class="ullinks">`
		} else {
			div += `<ul class="ullinks">`
		}

		for _, link := range set.Children {
			div += `<li class="ulchildlink">` + conversion.LinkAsAnchor(link)
			if link.Topic.Synopsis != "" {
				div += `<p>` + link.Topic.Synopsis + `</p>`
			}
			div += `</li>`
		}

		if set.CollType == dita.Sequence {
			div += `</ol>`
		} else {
			div += `</ul>`
		}
	}

	if hasFamilyLinks {
		div += `<div class="familylinks">`
		for _, set := range topic.Links {
			if set.Parent == nil && set.Prev == nil && set.Next == nil {
				continue
			}
			if set.Parent != nil {
				div += `<div class="parentlink"><strong>Parent topic: </strong>` + conversion.LinkAsAnchor(set.Parent) + `</div>`
			}
			if set.Prev != nil {
				div += `<div class="previouslink"><strong>Previous topic: </strong>` + conversion.LinkAsAnchor(set.Prev) + `</div>`
			}
			if set.Next != nil {
				div += `<div class="nextlink"><strong>Next topic: </strong>` + conversion.LinkAsAnchor(set.Next) + `</div>`
			}
		}
		div += `</div>`
	}

	grouped := make(map[string][]*ditaconvert.Link)
	order := []string{"video", "concept", "task", "reference", "information"}
	for _, set := range topic.Links {
		for _, link := range set.Siblings {
			kind := ""
			if link.Topic != nil && link.Topic.Original != nil {
				kind = link.Topic.Original.XMLName.Local
			}
			if link.Type != "" {
				kind = link.Type
			}
			if kind == "tutorial" {
				kind = "video"
			}
			if !contains(order, kind) {
				kind = "information"
			}

			grouped[kind] = append(grouped[kind], link)
		}
	}

	// for _, links := range grouped {
	// 	ditaconvert.SortLinks(links)
	// }

	for _, kind := range order {
		links := grouped[kind]
		if len(links) == 0 {
			continue
		}

		if kind != "information" {
			class := kindclass[kind]
			if len(links) > 1 {
				kind += "s"
			}
			div += `<div class="relinfo ` + class + `"><strong>Related ` + kind + `</strong>`
		} else {
			div += `<div class="relinfo"><strong>Related information</strong>`
		}
		for _, link := range links {
			div += "<div>" + conversion.LinkAsAnchor(link) + "</div>"
		}
		div += "</div>"
	}
	div += "</div>"

	return div
}

var kindclass = map[string]string{
	"video":  "reltutorials",
	"reference": "relref",
	"concept":   "relconcepts",
	"task":      "reltasks",
}

func (conversion *PageConversion) LinkAsAnchor(link *ditaconvert.Link) string {
	title := html.EscapeCharData(link.FinalTitle())
	if link.Scope == "external" {
		return `<a href="` + html.NormalizeURL(link.Href) + `" class="external-link" target="_blank" rel="nofollow">` + title + `</a>`
	}

	if link.Topic == nil {
		return `<span style="background: #f00">` + title + `</span>`
	}

	selector := link.Selector
	if selector != "" {
		selector = "#" + selector
	}

	slug, ok := conversion.Mapping.ByTopic[link.Topic]
	if !ok {
		return `<span style="background: #f00">` + title + `</span>`
	}

	return `<a href="` + string(slug) + selector + `" data-link="` + string(slug) + selector + `">` + title + `</a>`
}
