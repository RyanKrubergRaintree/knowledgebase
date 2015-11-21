package ditaconv

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/raintreeinc/knowledgebase/ditaconv/dita"
	"github.com/raintreeinc/knowledgebase/ditaconv/xmlconv"
	"github.com/raintreeinc/knowledgebase/extra/imagemap"
	"github.com/raintreeinc/knowledgebase/kb"
)

var alwaysHTML = map[string]bool{
	"div":       true,
	"code":      true,
	"pre":       true,
	"codeblock": true,
}

func (conv *convert) convertItem(decoder *xml.Decoder, start *xml.StartElement) {
	// NB! the converters must fully decode the element
	switch start.Name.Local {
	case "imagemap":
		_, err := conv.handleAttrs(start)
		conv.check(err)

		var m imagemap.XML
		err = decoder.DecodeElement(&m, start)
		conv.check(err)

		if err == nil {
			m.Image.Href = conv.convertImageURL(m.Image.Href)
			if item, err := imagemap.FromXML(&m); err == nil {
				conv.Page.Story.Append(item)
			} else {
				conv.check(err)
			}
		}
	case "data":
		switch datatype := getAttr(start, "datatype"); strings.ToLower(datatype) {
		case "rttutorial":
			conv.handleAttrs(start)
			href := getAttr(start, "href")
			conv.Page.Story.Append(kb.HTML("<video controls src=\"" + href + "\" >Browser doesn't support video.</video>"))
			xmlconv.Skip(decoder, start)
		default:
			text := conv.toHTML(decoder, start)
			conv.Page.Story.Append(kb.HTML(text))
			conv.Errors = append(conv.Errors, fmt.Errorf("unhandled datatype \"%v\"", datatype))
		}
	case "img", "image":
		_, err := conv.handleAttrs(start)
		conv.check(err)
		href := getAttr(start, "src")
		conv.Page.Story.Append(kb.Image("", href, ""))
		xmlconv.Skip(decoder, start)
	case "title":
		title, _ := xmlconv.Text(decoder, start)
		if title != "" {
			conv.Page.Story.Append(kb.HTML("<h3>" + title + "</h3>"))
		}
	case "xref", "link", "a":
		conv.handleAttrs(start)
		title, _ := xmlconv.Text(decoder, start)

		href := getAttr(start, "href")
		if title == "" {
			title = getAttr(start, "caption")
		}
		desc := getAttr(start, "title")
		conv.Page.Story.Append(kb.Reference(title, href, desc))
	default:
		text := conv.toHTML(decoder, start)
		if alwaysHTML[start.Name.Local] {
			conv.Page.Story.Append(kb.HTML(text))
		} else {
			// try to convert to paragraph
			if para, ok := asParagraph(text); ok {
				conv.Page.Story.Append(kb.Paragraph(para))
			} else {
				conv.Page.Story.Append(kb.HTML(text))
			}
		}
	}
}

func (conv *convert) toHTML(decoder *xml.Decoder, start *xml.StartElement) string {
	buf := bytes.Buffer{}
	enc := xmlconv.NewHTMLEncoder(&buf, conv.Rules)
	err := conv.Rules.ConvertElement(enc, decoder, start)
	conv.check(err)

	enc.Flush()
	return buf.String()
}

func (conv *convert) handleAttrs(start *xml.StartElement) (skip bool, err error) {
	href, title, desc, internal := "", "", "", false

	for i, a := range start.Attr {
		if a.Name.Local == "href" {
			if start.Name.Local == "img" || start.Name.Local == "image" || start.Name.Local == "fig" {
				start.Attr[i].Name.Local = "src"
				href = conv.convertImageURL(a.Value)
			} else {
				href, title, desc, internal = conv.convertLinkURL(a.Value)
			}
			start.Attr[i].Value = href
		} else if a.Name.Local == "id" {
			// id-s must be unique on a single web-page...
			// hence we convert it to "data-id", so that we can open same page
			// multiple times
			start.Attr[i].Name.Local = "data-id"
		}
	}

	if title != "" && getAttr(start, "caption") == "" {
		start.Attr = append(start.Attr, xml.Attr{xml.Name{Local: "caption"}, title})
	}

	if desc != "" && getAttr(start, "title") == "" {
		start.Attr = append(start.Attr, xml.Attr{xml.Name{Local: "title"}, desc})
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

func (conv *convert) convertLinkURL(url string) (href, title, desc string, internal bool) {
	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "mailto") {
		return url, "", "", false
	}

	var hash string
	i := strings.LastIndex(url, "#")
	if i >= 0 {
		url, hash = url[:i], url[i:]
	}

	var filename string
	if url != "" {
		filename = filepath.Join(filepath.Dir(conv.Topic.Filename), url)
	} else {
		filename = conv.Topic.Filename
	}

	if hash != "" {
		var err error
		full := conv.Index.Dir.Filepath(filename)
		title, err = dita.ExtractTitle(full, strings.TrimPrefix(hash, "#"))
		if err != nil {
			hash = ""
			conv.Errors = append(conv.Errors, fmt.Errorf("unable to extract title from %v [%v%v]: %s", filename, url, hash, err))
		}
	}

	topic, ok := conv.Mapping.Topics[canonicalName(filename)]
	if !ok {
		conv.Errors = append(conv.Errors, fmt.Errorf("did not find topic %v [%v%v]", filename, url, hash))
		return url + hash, title, "", false
	}

	if title == "" || hash == "" {
		title = topic.Title
		desc, _ = topic.Original.ShortDesc.Text()
	}

	slug, ok := conv.Mapping.ByTopic[topic]
	if !ok {
		conv.Errors = append(conv.Errors, fmt.Errorf("did not find topic %v [%v%v]", filename, url, hash))
		return url + hash, title, desc, false
	}

	return string(slug) + hash, title, desc, true
}
