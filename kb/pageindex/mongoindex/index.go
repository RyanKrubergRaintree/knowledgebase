package mongoindex

import (
	"errors"

	"github.com/egonelbre/fedwiki"
	"github.com/raintreeinc/knowledgebase/kb/pageindex"

	//TODO: migrate to v3
	"gopkg.in/mgo.v2-unstable"
	"gopkg.in/mgo.v2-unstable/bson"
)

// TODO: also implement pagestore
// var _ fedwiki.PageStore = (*Index)(nil)
var _ pageindex.Index = (*Index)(nil)

type Index struct {
	session    *mgo.Session
	collection string
}

func translate(err error) error {
	if err == nil {
		return nil
	}
	if err == mgo.ErrNotFound {
		return fedwiki.ErrNotExist
	}
	return err
}

// New returns a new MongoDB over a fedwiki/pagestore/mongostore
func New(url, collection string) (*Index, error) {
	main, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Index{main, collection}, nil
}

var fulltext = mgo.Index{
	Key: []string{
		"$text:title",
		"$text:synopsis",
		"$text:tags",
		"$text:story.text",
	},
	Weights: map[string]int{
		"title":      20,
		"tags":       10,
		"synopsis":   5,
		"story.text": 1,
	},
}

// Init creates the database indexes, if needed
func (index *Index) Init() error {
	session, c := index.c()
	defer session.Close()

	// TODO: add citations function
	return c.EnsureIndex(fulltext)
}

// returns the collection for storing items
func (index *Index) c() (*mgo.Session, *mgo.Collection) {
	s := index.session.Copy()
	return s, s.DB("").C(index.collection)
}

func (index *Index) Close() {
	index.session.Close()
}

func (index *Index) All() ([]*fedwiki.PageHeader, error) {
	session, c := index.c()
	defer session.Close()

	headers := []*fedwiki.PageHeader{}
	if err := c.Find(bson.M{}).Sort("_id").All(&headers); err != nil {
		return nil, translate(err)
	}

	return headers, nil
}

var taglisting = []bson.M{
	{"$project": bson.M{"tags": 1}},
	{"$unwind": "$tags"},
	{"$group": bson.M{
		"_id":   "$tags",
		"name":  bson.M{"$first": "$tags"},
		"count": bson.M{"$sum": 1}},
	},
	{"$match": bson.M{"count": bson.M{"$gt": 1}}},
	{"$sort": bson.M{"name": 1}},
}

func (index *Index) Tags() (tags []pageindex.TagInfo, err error) {
	session, c := index.c()
	defer session.Close()

	var taginfos []pageindex.TagInfo
	if err := c.Pipe(taglisting).All(&taginfos); err != nil {
		return nil, err
	}

	return taginfos, nil
}

func (index *Index) PagesByTag(tag string) ([]*fedwiki.PageHeader, error) {
	session, c := index.c()
	defer session.Close()

	headers := []*fedwiki.PageHeader{}
	if err := c.Find(bson.M{"meta.tags": tag}).All(&headers); err != nil {
		return nil, translate(err)
	}

	return headers, nil
}

func (index *Index) PagesCiting(slug fedwiki.Slug) ([]*fedwiki.PageHeader, error) {
	//session, c := index.c()
	//defer session.Close()
	//TODO:

	return nil, errors.New("Not implemented.")
}

func (index *Index) Search(content string) ([]*fedwiki.PageHeader, error) {
	session, c := index.c()
	defer session.Close()

	headers := []*fedwiki.PageHeader{}
	err := c.Find(bson.M{"$text": bson.M{"$search": content}}).
		Select(bson.M{"score": bson.M{"$meta": "textScore"}}).
		Sort("$textScore:score").
		All(&headers)
	if err != nil {
		return nil, translate(err)
	}
	return headers, nil
}

func (index *Index) RecentChanges(n int) ([]*fedwiki.PageHeader, error) {
	session, c := index.c()
	defer session.Close()

	headers := []*fedwiki.PageHeader{}
	if err := c.Find(bson.M{}).Sort("-modified").Limit(n).All(&headers); err != nil {
		return nil, translate(err)
	}

	return headers, nil
}
