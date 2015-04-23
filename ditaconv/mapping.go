package ditaconv

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/egonelbre/fedwiki"
	"github.com/raintreeinc/knowledgebase/ditaconv/xmlconv"
)

type Mapping struct {
	Index   *Index
	Topics  map[string]*Topic
	BySlug  map[fedwiki.Slug]*Topic
	ByTopic map[*Topic]fedwiki.Slug
	Rules   *xmlconv.Rules
}

func (m *Mapping) TopicsSorted() (r []*Topic) {
	for _, topic := range m.Topics {
		r = append(r, topic)
	}
	sort.Sort(byfilename(r))
	return r
}

type byfilename []*Topic

func (a byfilename) Len() int           { return len(a) }
func (a byfilename) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byfilename) Less(i, j int) bool { return a[i].Filename < a[j].Filename }

var (
	rxOr  = regexp.MustCompile(` ?/ ?`)
	rxAnd = regexp.MustCompile(`(?:[^\^]) ?& ?`)
)

// replace / and & inside the title
func titelize(title string) string {
	title = rxOr.ReplaceAllString(title, " or ")
	title = rxAnd.ReplaceAllString(title, " and ")
	return title
}

func CreateMapping(index *Index) (*Mapping, []error) {
	topics := index.Topics

	var errors []error
	byslug := make(map[fedwiki.Slug]*Topic, len(topics))
	bytopic := make(map[*Topic]fedwiki.Slug, len(topics))

	// assign slugs to the topics
	for _, topic := range topics {
		topic.Title = titelize(topic.Title)
		topic.ShortTitle = titelize(topic.ShortTitle)
		slug := fedwiki.Slugify(topic.Title)

		if other, clash := byslug[slug]; clash {
			errors = append(errors, fmt.Errorf("clashing title \"%v\" in \"%v\" and \"%v\"", topic.Title, topic.Filename, other.Filename))
			continue
		}

		if topic.Title == "" {
			errors = append(errors, fmt.Errorf("title missing in \"%v\"", topic.Filename))
			continue
		}

		byslug[slug] = topic
		bytopic[topic] = slug
	}

	// promote to shorter titles, if possible
	for prev, topic := range byslug {
		if topic.ShortTitle == "" || len(topic.Title) <= len(topic.ShortTitle) {
			continue
		}

		slug := fedwiki.Slugify(topic.ShortTitle)
		if _, exists := byslug[slug]; exists {
			continue
		}
		topic.Title = topic.ShortTitle
		topic.ShortTitle = ""

		delete(byslug, prev)
		byslug[slug] = topic
		bytopic[topic] = slug
	}

	m := &Mapping{
		Rules:   NewHTMLRules(),
		Index:   index,
		Topics:  topics,
		BySlug:  byslug,
		ByTopic: bytopic,
	}

	return m, errors
}
