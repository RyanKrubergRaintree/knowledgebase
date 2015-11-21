package ditaconv

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/raintreeinc/knowledgebase/ditaconv/xmlconv"
	"github.com/raintreeinc/knowledgebase/kb"
)

// Convert converts a topic to a federated wiki page
func (mapping *Mapping) Convert(topic *Topic) (page *kb.Page, fatal error, errs []error) {
	// make a shallow copy of rules
	rules := xmlconv.NewRules()
	rules.Merge(mapping.Rules)

	slug := mapping.ByTopic[topic]
	if slug[0] != '/' {
		slug = "/" + slug
	}
	convert := &convert{
		Page:    &kb.Page{},
		Slug:    slug,
		Topic:   topic,
		Index:   mapping.Index,
		Mapping: mapping,
		Rules:   rules,
	}

	convert.Rules.Handle.Element = convert.nestedElement
	convert.Rules.Callback["xref"] = convert.handleXRef

	convert.run()

	return convert.Page, convert.Fatal, convert.Errors
}

type convert struct {
	Page *kb.Page

	Slug    kb.Slug
	Topic   *Topic
	Index   *Index
	Mapping *Mapping

	Rules *xmlconv.Rules

	Errors []error
	Fatal  error
}

func (conv *convert) check(err error) bool {
	if err != nil {
		conv.Errors = append(conv.Errors, err)
		return true
	}
	return false
}

// entrypoint for starting the conversion
func (conv *convert) run() {
	info, err := conv.Index.Dir.Stat(conv.Topic.Filename)
	if err != nil {
		conv.Fatal = err
		return
	}

	topic := conv.Topic.Original
	if topic == nil {
		conv.Fatal = fmt.Errorf("no original topic")
		return
	}

	// find the body content
	bodytext := ""
	for _, node := range topic.Elements {
		if isBodyTag(node.XMLName.Local) {
			if bodytext != "" {
				conv.Errors = append(conv.Errors, fmt.Errorf("multiple body tags"))
				continue
			}
			bodytext = node.Content
		}
	}

	// create the page header
	conv.Page = &kb.Page{
		Slug:     conv.Slug,
		Title:    conv.Topic.Title,
		Modified: info.ModTime(),
		Synopsis: conv.Topic.Synopsis,
	}

	tags := convertTags(topic.Keywords)
	if len(tags) > 0 {
		conv.Page.Story.Append(kb.Tags(tags...))
	}

	if conv.Topic.Synopsis != "" {
		conv.Page.Story.Append(kb.Paragraph(conv.Topic.Synopsis))
	}

	defer func() {
		if x := recover(); x != nil {
			conv.Fatal = fmt.Errorf("fatal: %v", x)
		}
	}()

	conv.parse(bodytext)
	conv.addRelatedLinks()
	conv.dropEmpty()
	conv.assignIDs()
}

// splits node recursively into multiple story items
func (conv *convert) parse(text string) {
	decoder := xml.NewDecoder(strings.NewReader(text))
	conv.unwrap(decoder, nil)
}

// splits node recursively into multiple story items
func (conv *convert) unwrap(decoder *xml.Decoder, start *xml.StartElement) {
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return
		}
		if conv.check(err) {
			return
		}

		switch token := token.(type) {
		case xml.StartElement:
			if shouldUnwrap(token.Name) || conv.Rules.Unwrap[token.Name.Local] {
				conv.unwrap(decoder, &token)
			} else {
				conv.convertItem(decoder, &token)
			}
		case xml.CharData, xml.Comment, xml.ProcInst, xml.Directive:
			// ignore
		case xml.EndElement:
			return
		}
	}
}

func (conv *convert) addRelatedLinks() {
	for _, set := range conv.Topic.Links {
		text := "<h4>Related</h4>"
		text += "<ul>"

		links := []*Topic{set.Parent, set.Prev, set.Next}
		links = append(links, set.Children...)
		links = append(links, set.Siblings...)

		for _, topic := range links {
			if topic == nil {
				continue
			}
			text += "<li>" + conv.asLink(topic) + "</li>"
		}
		text += "</ul>"

		conv.Page.Story.Append(kb.HTML(text))
	}
}

func (conv *convert) asLink(topic *Topic) string {
	slug := string(conv.Mapping.ByTopic[topic])
	title := topic.Title
	return "<a href=\"" + slug + "\" data-link=\"" + slug + "\">" + title + "</a>"
}

func (conv *convert) dropEmpty() {
	s := conv.Page.Story
	s = s[:0:cap(s)]
	for _, item := range conv.Page.Story {
		if item.Type() == "paragraph" || item.Type() == "html" {
			if item["text"] != "" {
				s.Append(item)

			}
		} else {
			s.Append(item)
		}
	}
	conv.Page.Story = s
}

func (conv *convert) assignIDs() {
	s := conv.Page.Story
	for _, item := range s {
		item["id"] = kb.NewID()
	}
}
