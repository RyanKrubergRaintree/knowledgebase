package dita

import "encoding/xml"

// MapNode represents an xml element inside a .ditamap
type MapNode struct {
	XMLName  xml.Name
	Title    string `xml:"title"`
	NavTitle string `xml:"navtitle,attr"`
	Href     string `xml:"href,attr"`

	Type     string         `xml:"type,attr"`
	CollType CollectionType `xml:"collection-type,attr"`
	Linking  Linking        `xml:"linking,attr"`

	Format    string `xml:"format,attr"`
	TOC       string `xml:"toc,attr"`
	LockTitle string `xml:"locktitle,attr"`

	Children []*MapNode `xml:",any"`
}

var (
	RelTable   = xml.Name{Local: "reltable"}
	TopicRef   = xml.Name{Local: "topicref"}
	MapRef     = xml.Name{Local: "mapref"}
	TopicHead  = xml.Name{Local: "topichead"}
	TopicGroup = xml.Name{Local: "topicgroup"}
)

type CollectionType string

const (
	// parent <-> child
	Unordered = CollectionType("unordered")
	// parent <-> child
	// child <-> child
	Family = CollectionType("family")
	// parent <-> child
	// child <-> child + 1
	Sequence = CollectionType("sequence")
)

type Linking string

func (linking Linking) CanLinkFrom() bool {
	return linking != NoLinking && linking != TargetOnly
}

func (linking Linking) CanLinkTo() bool {
	return linking != NoLinking && linking != SourceOnly
}

const (
	NormalLinking = Linking("normal")
	NoLinking     = Linking("none")
	SourceOnly    = Linking("sourceonly")
	TargetOnly    = Linking("targetonly")
)
