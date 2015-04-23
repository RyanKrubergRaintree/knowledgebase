package ditaconv

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"path"
	"strings"

	"github.com/egonelbre/fedwiki/item"

	"github.com/raintreeinc/knowledgebase/ditaconv/xmlconv"
)

var alwaysHTML = map[string]bool{
	"div":       true,
	"code":      true,
	"pre":       true,
	"codeblock": true,
}

func (conv *convert) convertItem(decoder *xml.Decoder, start *xml.StartElement) {
	switch start.Name.Local {
	case "data":
		switch datatype := getAttr(start, "datatype"); strings.ToLower(datatype) {
		case "rttutorial":
			conv.handleAttrs(start)
			href := getAttr(start, "href")
			conv.Page.Story.Append(item.HTML("<video controls src=\"" + href + "\" >Browser doesn't support video.</video>"))
		default:
			text := conv.toHTML(decoder, start)
			conv.Page.Story.Append(item.HTML(text))
			conv.Errors = append(conv.Errors, fmt.Errorf("unhandled datatype \"%v\"", datatype))
		}
	case "img", "image":
		_, err := conv.handleAttrs(start)
		conv.check(err)
		href := getAttr(start, "src")
		conv.Page.Story.Append(item.Image("", href, ""))
	case "title":
		title, _ := xmlconv.Text(decoder, start)
		if title != "" {
			conv.Page.Story.Append(item.HTML("<h3>" + title + "</h3>"))
		}
	case "xref", "link", "a":
		title, _ := xmlconv.Text(decoder, start)
		if title == "" {
			title = getAttr(start, "href")
		}
		href := getAttr(start, "href")
		conv.Page.Story.Append(item.Reference(title, href, ""))
	default:
		text := conv.toHTML(decoder, start)
		if alwaysHTML[start.Name.Local] {
			conv.Page.Story.Append(item.HTML(text))
		} else {
			// try to convert to paragraph
			if para, ok := asParagraph(text); ok {
				conv.Page.Story.Append(item.Paragraph(para))
			} else {
				conv.Page.Story.Append(item.HTML(text))
			}
		}
	}
}

func (conv *convert) toHTML(decoder *xml.Decoder, start *xml.StartElement) string {
	buf := bytes.Buffer{}
	enc := xmlconv.NewHTMLEncoder(&buf)
	err := conv.Rules.ConvertElement(enc, decoder, start)
	conv.check(err)

	enc.Flush()
	return buf.String()
}

func (conv *convert) handleAttrs(start *xml.StartElement) (skip bool, err error) {
	internal := false
	href := ""
	for i, a := range start.Attr {
		if a.Name.Local == "href" {
			if start.Name.Local == "img" || start.Name.Local == "image" || start.Name.Local == "fig" {
				start.Attr[i].Name.Local = "src"
				href = conv.convertImageURL(a.Value)
			} else {
				href, internal = conv.convertLinkURL(a.Value)
			}
			start.Attr[i].Value = href
		}
	}

	if internal {
		start.Attr = append(start.Attr, xml.Attr{xml.Name{Local: "data-link"}, href})
	}

	start.Attr = append(start.Attr, xml.Attr{xml.Name{Local: "data-dita"}, start.Name.Local})

	return false, nil
}

func (conv *convert) convertImageURL(url string) (href string) {
	// if it's a remote link then preserve it
	if strings.HasPrefix(url, "http") {
		return url
	}

	name := path.Join(path.Dir(conv.Topic.Filename), url)
	data, err := conv.Index.Dir.ReadFile(name)
	if err != nil {
		conv.check(err)
		return path.Clean(name)
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	ext := strings.ToLower(path.Ext(name))
	if ext == "" {
		conv.Errors = append(conv.Errors, fmt.Errorf("invalid image link"))
		return path.Clean(name)
	}

	ext = ext[1:]
	switch ext {
	case "jpg", "jpeg":
		return "data:image/jpeg;base64," + encoded
	default:
		return "data:image/" + ext + ";base64," + encoded
	}
}

func (conv *convert) convertLinkURL(url string) (href string, internal bool) {
	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "mailto") {
		return url, false
	}

	var hash string
	i := strings.LastIndex(url, "#")
	if i >= 0 {
		url, hash = url[:i], url[i:]
	}

	if url == "" {
		//TODO: implement internal reference links
		return url + hash, false
	}

	name := path.Clean(path.Join(path.Dir(conv.Topic.Filename), url))
	cname := strings.ToLower(name)

	topic, ok := conv.Mapping.Topics[cname]
	if !ok {
		conv.Errors = append(conv.Errors, fmt.Errorf("did not find topic %v [%v%v]", name, url, hash))
		return url + hash, false
	}

	slug, ok := conv.Mapping.ByTopic[topic]
	if !ok {
		conv.Errors = append(conv.Errors, fmt.Errorf("did not find topic %v [%v%v]", name, url, hash))
		return url + hash, false
	}

	return string(slug) + hash, true
}
