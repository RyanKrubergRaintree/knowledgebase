package ditaconv

import (
	"github.com/raintreeinc/knowledgebase/ditaconv/dita"
	"github.com/raintreeinc/knowledgebase/kb"
)

type Links struct {
	CollType   dita.CollectionType
	Parent     *Topic
	Prev, Next *Topic
	Siblings   []*Topic
	Children   []*Topic
}

func (links *Links) IsEmpty() bool {
	return links.Parent == nil &&
		links.Prev == nil && links.Next == nil &&
		len(links.Siblings) == 0 &&
		len(links.Children) == 0
}

func CreateLinks(context Context, entries []*Entry) {
	linkable := make([]*Entry, 0, len(entries))
	for _, e := range entries {
		if e.Topic != nil {
			linkable = append(linkable, e)
		}
	}

	if context.Entry.Topic != nil && context.Entry.Linking.CanLinkTo() {
		for _, a := range linkable {
			if a.Linking.CanLinkFrom() {
				a.Topic.Links = append(a.Topic.Links, Links{Parent: context.Entry.Topic})
			}
		}
	}

	if context.Entry.Topic != nil && context.Entry.Linking.CanLinkFrom() {
		links := Links{CollType: context.CollType}
		for _, a := range linkable {
			if a.Linking.CanLinkTo() {
				links.Children = append(links.Children, a.Topic)
			}
		}

		if len(links.Children) > 0 {
			context.Entry.Topic.Links = append(context.Entry.Topic.Links, links)
		}
	}

	switch context.CollType {
	case dita.Family:
		for _, a := range linkable {
			links := Links{}
			if !a.Linking.CanLinkFrom() {
				continue
			}

			for _, b := range linkable {
				if a != b && b.Linking.CanLinkTo() {
					links.Siblings = append(links.Siblings, b.Topic)
				}
			}

			if len(links.Siblings) > 0 {
				a.Topic.Links = append(a.Topic.Links, links)
			}
		}
	case dita.Sequence:
		for i, a := range linkable {
			if !a.Linking.CanLinkFrom() {
				continue
			}

			links := Links{}
			if i-1 >= 0 {
				prev := linkable[i-1]
				if prev.Linking.CanLinkTo() {
					links.Prev = prev.Topic
				}
			}
			if i+1 < len(linkable) {
				next := linkable[i+1]
				if next.Linking.CanLinkTo() {
					links.Next = next.Topic
				}
			}
			if links.Prev != nil || links.Next != nil {
				a.Topic.Links = append(a.Topic.Links, links)
			}
		}
	}
}

func InterLink(A, B []*Entry) {
	for _, a := range A {
		if a.Topic == nil || !a.Linking.CanLinkFrom() {
			continue
		}

		links := Links{}
		for _, b := range B {
			if b.Topic == nil {
				continue
			}

			if b.Linking.CanLinkTo() {
				links.Siblings = append(links.Siblings, b.Topic)
			}
		}

		if len(links.Siblings) > 0 {
			a.Topic.Links = append(a.Topic.Links, links)
		}
	}
}

func emptylinks(links []Links) bool {
	for _, set := range links {
		if !set.IsEmpty() {
			return false
		}
	}
	return true
}

func (conv *convert) addRelatedLinks() {
	if emptylinks(conv.Topic.Links) {
		return
	}

	text := "<div class=\"dita-related-links\">"

	for _, set := range conv.Topic.Links {
		if len(set.Children) > 0 {
			if set.CollType == dita.Sequence {
				text += "<ol>"
			} else {
				text += "<ul>"
			}

			for _, topic := range set.Children {
				text += "<li>" + conv.asLink(topic)
				if topic.Synopsis != "" {
					text += "<p>" + topic.Synopsis + "</p>"
				}
				text += "</li>"
			}
			text += "</ol>"

			if set.CollType == dita.Sequence {
				text += "</ol>"
			} else {
				text += "</ul>"
			}
		}

		if set.Parent != nil {
			text += "<div><b>Parent topic: </b>" + conv.asLink(set.Parent) + "</div>"
		}
		if set.Prev != nil {
			text += "<div><b>Previous topic: </b>" + conv.asLink(set.Prev) + "</div>"
		}
		if set.Next != nil {
			text += "<div><b>Next topic: </b>" + conv.asLink(set.Next) + "</div>"
		}
	}

	grouped := make(map[string][]*Topic)
	order := []string{"concept", "task", "reference", "tutorial", "information"}
	for _, set := range conv.Topic.Links {
		for _, topic := range set.Siblings {
			kind := topic.Original.XMLName.Local
			if !contains(order, kind) {
				kind = "information"
			}

			grouped[kind] = append(grouped[kind], topic)
		}
	}

	for _, kind := range order {
		topics := grouped[kind]
		if len(topics) == 0 {
			continue
		}

		if kind != "information" {
			kind += "s"
		}
		text += "<div><b>Related " + kind + "</b>"
		for _, topic := range topics {
			text += "<div>" + conv.asLink(topic) + "</div>"
		}
		text += "</div>"
	}

	text += "</div>"
	conv.Page.Story.Append(kb.HTML(text))
}

func (conv *convert) asLink(topic *Topic) string {
	slug := string(conv.Mapping.ByTopic[topic])
	title := topic.Title
	return "<a href=\"" + slug + "\" data-link=\"" + slug + "\">" + title + "</a>"
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
