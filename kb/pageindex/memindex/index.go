package memindex

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/egonelbre/fedwiki"

	"github.com/raintreeinc/knowledgebase/kb/pageindex"
)

var _ pageindex.Index = (*Index)(nil)

type Index struct {
	store fedwiki.PageStore
	tick  *time.Ticker

	info atomic.Value
}

func New(store fedwiki.PageStore, updateInterval time.Duration) *Index {
	index := &Index{store: store, tick: time.NewTicker(updateInterval)}
	go index.reloader()
	return index
}

func (index *Index) Close() { index.tick.Stop() }

func (index *Index) reloader() {
	index.reload()
	for range index.tick.C {
		index.reload()
	}
}

func (index *Index) reload() {
	index.info.Store(newInfo(index.store))
}

func (index *Index) getInfo() *info {
	return index.info.Load().(*info)
}

func (index *Index) All() ([]*fedwiki.PageHeader, error) {
	info := index.getInfo()
	if info == nil {
		return nil, errors.New("Index not loaded.")
	}

	return info.bySlug, nil
}

func (index *Index) Tags() (tags []pageindex.TagInfo, err error) {
	info := index.getInfo()
	if info == nil {
		return nil, errors.New("Index not loaded.")
	}

	for tagname, pages := range info.tags {
		tags = append(tags, pageindex.TagInfo{tagname, len(pages)})
	}
	return tags, nil
}

func (index *Index) PagesByTag(tag string) ([]*fedwiki.PageHeader, error) {
	info := index.getInfo()
	if info == nil {
		return nil, errors.New("Index not loaded.")
	}

	slugs, ok := info.tags[tag]
	if !ok {
		return nil, errors.New("Tag does not exist.")
	}

	headers := make([]*fedwiki.PageHeader, len(slugs))
	for _, slug := range slugs {
		header := info.header[slug]
		headers = append(headers, &header)
	}
	return headers, nil
}

func (index *Index) PagesCiting(slug fedwiki.Slug) ([]*fedwiki.PageHeader, error) {
	info := index.getInfo()
	if info == nil {
		return nil, errors.New("Index not loaded.")
	}

	slugs, ok := info.citations[slug]
	if !ok {
		return nil, nil
	}

	headers := make([]*fedwiki.PageHeader, len(slugs))
	for _, slug := range slugs {
		header := info.header[slug]
		headers = append(headers, &header)
	}
	return headers, nil
}

func (index *Index) Search(content string) ([]*fedwiki.PageHeader, error) {
	info := index.getInfo()
	if info == nil {
		return nil, errors.New("Index not loaded.")
	}

	//TODO: implement
	return nil, nil
}

func (index *Index) RecentChanges(n int) ([]*fedwiki.PageHeader, error) {
	info := index.getInfo()
	if info == nil {
		return nil, errors.New("Index not loaded.")
	}

	if n > len(info.byDate) {
		n = len(info.byDate)
	}

	return info.byDate[:n], nil
}
