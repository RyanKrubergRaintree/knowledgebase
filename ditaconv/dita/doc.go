package dita

import (
	"encoding/xml"
	"io/ioutil"
	"log"
)

type Topic struct {
	Title     string   `xml:"title"`
	NavTitle  string   `xml:"titlealts>navtitle"`
	Keywords  []string `xml:"prolog>metadata>keywords>indexterm"`
	ShortDesc string   `xml:"shortdesc"`

	RelatedLink []Link `xml:"related-links>link>href"`

	Elements []Body `xml:",any"`
}

type TopicHeader struct {
	Title     string `xml:"title"`
	ShortDesc string `xml:"shortdesc"`
	TitleAlts struct {
		Nav    string `xml:"navtitle"`
		Search string `xml:"searchtitle"`
	} `xml:"titlealts"`
}

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
