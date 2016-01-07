package dita

import (
	"fmt"
	"sort"

	"github.com/raintreeinc/ditaconvert"
	"github.com/raintreeinc/knowledgebase/kb"
)

type TitleMapping struct {
	Topics  map[string]*ditaconvert.Topic
	BySlug  map[kb.Slug]*ditaconvert.Topic
	ByTopic map[*ditaconvert.Topic]kb.Slug
}

func NewTitleMapping() *TitleMapping {
	return &TitleMapping{
		Topics:  make(map[string]*ditaconvert.Topic),
		BySlug:  make(map[kb.Slug]*ditaconvert.Topic),
		ByTopic: make(map[*ditaconvert.Topic]kb.Slug),
	}
}

type byTopicPath []*ditaconvert.Topic

func (a byTopicPath) Len() int           { return len(a) }
func (a byTopicPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTopicPath) Less(i, j int) bool { return a[i].Path < a[j].Path }

func (m *TitleMapping) TopicsSorted() (r []*ditaconvert.Topic) {
	for _, topic := range m.Topics {
		r = append(r, topic)
	}
	sort.Sort(byTopicPath(r))
	return r
}

func RemapTitles(conversion *Conversion, index *ditaconvert.Index) (*TitleMapping, []error) {
	var errors []error

	mapping := NewTitleMapping()

	// assign slugs to topics
	for _, topic := range index.Topics {
		slug := kb.Slugify(topic.Title)
		if other, clash := mapping.BySlug[slug]; clash {
			errors = append(errors, fmt.Errorf("clashing title \"%v\" in \"%v\" and \"%v\"", topic.Title, topic.Path, other.Path))
			continue
		}

		if topic.Title == "" {
			errors = append(errors, fmt.Errorf("title missing in \"%v\"", topic.Filename))
			continue
		}

		mapping.BySlug[slug] = topic
		mapping.ByTopic[topic] = slug
	}

	// promote to shorter titles, if possible
	for prev, topic := range mapping.BySlug {
		if topic.ShortTitle == "" || len(topic.Title) <= len(topic.ShortTitle) {
			continue
		}

		slug := kb.Slugify(topic.ShortTitle)
		if _, exists := mapping.BySlug[slug]; exists {
			continue
		}
		topic.Title = topic.ShortTitle
		topic.ShortTitle = ""

		delete(mapping.BySlug, prev)
		mapping.BySlug[slug] = topic
		mapping.byTopic[topic] = slug
	}

	return mapping, errors
}
