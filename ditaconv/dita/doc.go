package dita

import (
	"encoding/xml"
	"io/ioutil"
	"log"

	"github.com/raintreeinc/knowledgebase/ditaconv/xmlconv"
)

type Topic struct {
	XMLName   xml.Name
	Title     string   `xml:"title"`
	NavTitle  string   `xml:"titlealts>navtitle"`
	Keywords  []string `xml:"prolog>metadata>keywords>indexterm"`
	ShortDesc InnerXML `xml:"shortdesc"`

	RelatedLink []Link `xml:"related-links>link"`

	Elements []Body `xml:",any"`
	Raw      []byte
}

type TopicHeader struct {
	Title     string   `xml:"title"`
	ShortDesc InnerXML `xml:"shortdesc"`
	TitleAlts struct {
		Nav    string `xml:"navtitle"`
		Search string `xml:"searchtitle"`
	} `xml:"titlealts"`
}

type InnerXML struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

func (x *InnerXML) Text() (string, error) { return xmlconv.StripTags(x.Content) }

type Body struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

type Link struct {
	Href   string `xml:"href,attr"`
	Format string `xml:"format,attr,omitempty"`
	Scope  string `xml:"scope,attr,omitempty"`
	Text   string `xml:"linktext,omitempty"`
}

func LoadTopic(filename string) (*Topic, error) {
	topic := &Topic{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	topic.Raw = data

	data, err = replaceConrefs(data, filename)
	if err != nil {
		log.Println(err)
	}

	err = xml.Unmarshal(data, topic)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

func LoadTopicHeader(filename string) (*TopicHeader, error) {
	header := &TopicHeader{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, header)
	if err != nil {
		return nil, err
	}

	return header, nil
}
