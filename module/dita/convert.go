package dita

import (
	"encoding/xml"
	"fmt"
	"html"
	"path"
	"strings"

	"github.com/raintreeinc/ditaconvert"
	"github.com/raintreeinc/knowledgebase/kb"
)

type PageConversion struct {
	*Conversion
	Mapping *TitleMapping
	Slug    kb.Slug
	Index   *ditaconvert.Index
	Topic   *ditaconvert.Topic
	Context *ditaconvert.Context
}

func (conversion *PageConversion) Convert() (page *kb.Page, errs []error, fatal error) {
	conversion.Context = ditaconvert.NewConversion(conversion.Index, conversion.Topic)
	conversion.Context.Encoder.RewriteID = "data-id"

	context, topic := conversion.Context, conversion.Topic

	page = &kb.Page{
		Slug:     conversion.Slug,
		Title:    topic.Title,
		Modified: topic.Modified,
		Synopsis: topic.Synopsis,
	}

	context.Rules.Custom["a"] = conversion.ToSlug
	context.Rules.Custom["img"] = conversion.InlineImage
	context.Rules.Custom["imagemap"] = conversion.ConvertImageMap

	if err := context.Run(); err != nil {
		return page, nil, err
	}

	if tags := conversion.ConvertTags(); len(tags) > 0 {
		page.Story.Append(kb.Tags(tags...))
	}

	page.Story.Append(kb.HTML(context.Output.String()))
	page.Story.Append(kb.HTML(conversion.RelatedLinksAsHTML()))

	return page, context.Errors, nil
}

func (conversion *PageConversion) ConvertTags() []string {
	raw := conversion.Topic.Original.Prolog.Keywords.Terms()
	for _, key := range conversion.Topic.Original.Prolog.ResourceID {
		raw = append(raw, "id/"+key.Name)
	}

	tags := []string{}
	for _, tag := range raw {
		slug := string(kb.Slugify(tag))
		if slug == "" || slug == "-" {
			continue
		}
		tags = append(tags, tag)
	}

	return tags
}

func (conversion *PageConversion) ToSlug(context *ditaconvert.Context, dec *xml.Decoder, start xml.StartElement) error {
	var href, title, desc string
	var internal bool

	href = getAttr(&start, "href")
	if href != "" {
		href, title, desc, internal = conversion.ResolveLinkInfo(href)
		setAttr(&start, "href", href)
	}

	if desc != "" && getAttr(&start, "title") == "" {
		setAttr(&start, "title", desc)
	}

	setAttr(&start, "scope", "")
	if internal && href != "" {
		setAttr(&start, "data-link", href)
	}

	if !internal {
		if class := getAttr(&start, "class"); class != "" {
			setAttr(&start, "class", class+" external-link")
		} else {
			setAttr(&start, "class", "external-link")
		}
	}

	if getAttr(&start, "format") != "" && href != "" {
		setAttr(&start, "format", "")
		ext := strings.ToLower(path.Ext(href))
		if ext == ".doc" || ext == ".xml" || ext == ".rtf" || ext == ".zip" || ext == ".exe" {
			setAttr(&start, "download", path.Base(href))
		} else {
			setAttr(&start, "target", "_blank")
		}
	}
	// encode starting tag and attributes
	if err := context.Encoder.WriteStart("a", start.Attr...); err != nil {
		return err
	}

	// recurse on child tokens
	err, count := context.RecurseChildCount(dec)
	if err != nil {
		return err
	}
	if count == 0 {
		context.Encoder.WriteRaw(html.EscapeString(title))
	}
	return context.Encoder.WriteEnd("a")
}

func (conversion *PageConversion) InlineImage(context *ditaconvert.Context, dec *xml.Decoder, start xml.StartElement) error {
	href := getAttr(&start, "href")
	setAttr(&start, "src", context.InlinedImageURL(href))
	setAttr(&start, "href", "")

	placement := getAttr(&start, "placement")
	setAttr(&start, "placement", "")
	if placement == "break" {
		context.Encoder.WriteStart("p",
			xml.Attr{Name: xml.Name{Local: "class"}, Value: "image"})
	}

	err := context.EmitWithChildren(dec, start)

	if placement == "break" {
		context.Encoder.WriteEnd("p")
	}

	return err
}

func (conversion *PageConversion) ResolveLinkInfo(url string) (href, title, synopsis string, internal bool) {
	if strings.HasPrefix(url, "http:") || strings.HasPrefix(url, "https:") || strings.HasPrefix(url, "mailto:") {
		return url, "", "", false
	}
	context := conversion.Context

	var selector, hash string
	url, selector = ditaconvert.SplitLink(url)
	if selector != "" {
		hash = "#" + selector
	}

	if url == "" {
		return hash, "", "", true
	}

	name := context.DecodingPath
	if url != "" {
		name = path.Join(path.Dir(context.DecodingPath), url)
	}

	topic, ok := context.Index.Topics[ditaconvert.CanonicalPath(name)]
	if !ok {
		context.Errors = append(context.Errors,
			fmt.Errorf("did not find topic %v [%v%v]", name, url, selector))
		return "", "", "", false
	}

	if selector != "" {
		var err error
		title, err = ditaconvert.ExtractTitle(topic.Raw, selector)
		if err != nil {
			context.Errors = append(context.Errors,
				fmt.Errorf("unable to extract title from %v [%v%v]: %v", name, url, selector, err))
		}
	}

	if title == "" && topic.Original != nil {
		title = topic.Title
		if selector == "" {
			synopsis, _ = topic.Original.ShortDesc.Text()
		}
	}

	slug, ok := conversion.Mapping.ByTopic[topic]
	if !ok {
		return href, title, synopsis, false
	}

	return string(slug) + hash, title, synopsis, true
}
