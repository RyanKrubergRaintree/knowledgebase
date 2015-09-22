package ditaconv

import (
	"net/url"
	"path"
	"path/filepath"

	"github.com/raintreeinc/knowledgebase/ditaconv/dita"
)

type Topic struct {
	Filename string

	Title      string
	ShortTitle string
	Synopsis   string

	Links []Links

	Original *dita.Topic
}

type Links struct {
	Parent     *Topic
	Prev, Next *Topic
	Siblings   []*Topic
	Children   []*Topic
}

type Context struct {
	Dir     string
	Entry   *Entry
	Type    dita.CollectionType
	Linking dita.Linking
	TOC     bool
}

type Map struct {
	Filename string
	Entries  []*Entry
	Node     *dita.MapNode
}

type Index struct {
	Dir  Dir
	Root string

	Nav *Entry

	// path --> entry
	Maps   map[string]*Map
	Topics map[string]*Topic

	Errors []error
}

type Entry struct {
	Title     string
	Topic     *Topic
	Linking   dita.Linking
	TOC       bool
	LockTitle bool

	Children []*Entry
}

func (index *Index) check(err error) bool {
	if err != nil {
		index.Errors = append(index.Errors, err)
		return true
	}
	return false
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
		links := Links{}
		for _, a := range linkable {
			if a.Linking.CanLinkTo() {
				links.Children = append(links.Children, a.Topic)
			}
		}

		if len(links.Children) > 0 {
			context.Entry.Topic.Links = append(context.Entry.Topic.Links, links)
		}
	}

	switch context.Type {
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

func (index *Index) processRelRow(context Context, node *dita.MapNode) {
	var entrysets [][]*Entry
	for _, cell := range node.Children {
		if cell.XMLName.Local != "relcell" {
			continue
		}

		var entries []*Entry
		for _, child := range cell.Children {
			subentries := index.processNode(context, child)
			entries = append(entries, subentries...)
		}

		entrysets = append(entrysets, entries)
	}

	for i, a := range entrysets {
		for j, b := range entrysets {
			if i != j {
				InterLink(a, b)
			}
		}
	}
}

func (index *Index) processNode(context Context, node *dita.MapNode) []*Entry {
	if node == nil {
		panic("shouldn't happen")
	}

	if node.Format != "" {
		return []*Entry{}
	}

	href, err := url.QueryUnescape(node.Href)
	if !index.check(err) {
		node.Href = href
	}

	if node.Type != "" {
		context.Type = node.Type
	} else {
		context.Type = dita.Unordered
	}

	if node.Linking != "" {
		context.Linking = node.Linking
	}

	if node.XMLName == dita.TopicGroup || node.XMLName.Local == "map" {
		var entries []*Entry
		childcontext := context
		childcontext.TOC = isChildTOC(context.TOC, node.TOC)
		for _, child := range node.Children {
			entries = append(entries, index.processNode(childcontext, child)...)
		}
		context.Entry = &Entry{}
		CreateLinks(context, entries)
		return entries
	}

	if node.XMLName == dita.MapRef {
		childcontext := context
		childcontext.TOC = isChildTOC(context.TOC, node.TOC)
		return index.loadMap(childcontext, node.Href)
	}

	if node.XMLName == dita.RelTable {
		for _, row := range node.Children {
			if row.XMLName.Local == "relrow" {
				index.processRelRow(context, row)
			}
		}
		return []*Entry{}
	}

	entry := &Entry{
		Title:     node.NavTitle,
		LockTitle: node.LockTitle == "yes",
		TOC:       isChildTOC(context.TOC, node.TOC),
	}
	if entry.Title == "" {
		entry.Title = node.Title
	}
	entry.Linking = context.Linking

	if node.Href != "" {
		entry.Topic = index.loadTopic(context, node.Href)
		if entry.Title == "" && entry.Topic != nil && !entry.LockTitle {
			entry.Title = entry.Topic.Title
		}
	}

	childcontext := context
	childcontext.Entry = entry
	childcontext.TOC = entry.TOC
	var entries []*Entry
	for _, child := range node.Children {
		subentries := index.processNode(childcontext, child)
		entries = append(entries, subentries...)
	}
	entry.Children = append(entry.Children, entries...)

	CreateLinks(childcontext, entries)

	return []*Entry{entry}
}

func (index *Index) loadMap(context Context, filename string) []*Entry {
	name := path.Join(context.Dir, filename)
	cname := canonicalName(name)
	if m, loaded := index.Maps[cname]; loaded {
		return m.Entries
	}

	m, err := index.Dir.LoadMap(name)
	if index.check(err) {
		return nil
	}

	index.Maps[cname] = m

	context.Dir = path.Dir(name)
	m.Entries = index.processNode(context, m.Node)

	return m.Entries
}

// Loads a single topic from a concrete file with context
func (index *Index) loadTopic(context Context, filename string) *Topic {
	name := path.Join(context.Dir, filename)
	cname := canonicalName(name)
	if topic, loaded := index.Topics[cname]; loaded {
		return topic
	}

	topic, err := index.Dir.LoadTopic(name)
	index.check(err)
	index.Topics[cname] = topic

	return topic
}

// Load loads the full index and linked maps starting from "filename"
func LoadIndex(filename string) (*Index, []error) {
	index := &Index{
		Dir:  Dir(filepath.Dir(filename)),
		Root: filepath.Base(filename),

		Nav: &Entry{
			Title:   "Navigation",
			Linking: dita.NormalLinking,
		},

		Maps:   make(map[string]*Map),
		Topics: make(map[string]*Topic),
	}

	context := Context{
		Dir:     "",
		Entry:   index.Nav,
		Type:    dita.Unordered,
		Linking: dita.NormalLinking,
		TOC:     true,
	}

	entries := index.loadMap(context, filepath.Base(filename))
	index.Nav.Children = append(index.Nav.Children, entries...)

	return index, index.Errors
}

func isChildTOC(parenttoc bool, childtoc string) bool {
	if childtoc == "" {
		return parenttoc
	}
	return childtoc == "yes"
}
