package dita

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/bradfitz/slice"
	"github.com/raintreeinc/ditaconvert"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kb/items/index"
)

type Conversion struct {
	Group   kb.Slug
	Ditamap string

	Pages map[kb.Slug]*kb.Page
	Raw   map[kb.Slug][]byte
	Slugs []kb.Slug
	Nav   *index.Item

	LoadErrors    []error
	MappingErrors []error
	Errors        []ConversionError
}

func NewConversion(group kb.Slug, ditamap string) *Conversion {
	return &Conversion{
		Group:   group,
		Ditamap: ditamap,
		Pages:   make(map[kb.Slug]*kb.Page),
		Raw:     make(map[kb.Slug][]byte),
	}
}

type ConversionError struct {
	Path   string
	Slug   kb.Slug
	Fatal  error
	Errors []error
}

const maxPageSize = 1048575

func (context *Conversion) Run() {
	fs := ditaconvert.Dir(filepath.Dir(context.Ditamap))
	index := ditaconvert.NewIndex(fs)
	index.LoadMap(filepath.Base(context.Ditamap))

	context.LoadErrors = index.Errors

	mapping, mappingErrors := RemapTitles(context, index)
	context.MappingErrors = mappingErrors

	for slug, topic := range mapping.BySlug {
		page, errs, fatal := (&PageConversion{
			Conversion: context,
			Mapping:    mapping,
			Slug:       slug,
			Index:      index,
			Topic:      topic,
		}).Convert()

		if fatal != nil {
			context.Errors = append(context.Errors, ConversionError{
				Path:  topic.Path,
				Slug:  slug,
				Fatal: fatal,
			})
		} else if len(errs) > 0 {
			context.Errors = append(context.Errors, ConversionError{
				Path:   topic.Path,
				Slug:   slug,
				Errors: errs,
			})
		}

		data, err := json.Marshal(page)
		if err != nil {
			log.Println(err)
		}

		if len(data) > maxPageSize {
			context.Errors = append(context.Errors, ConversionError{
				Path:  topic.Path,
				Slug:  slug,
				Fatal: fmt.Errorf("Page is too large %vMB (%vB)", len(data)>>20, len(data)),
			})
		}

		context.Pages[slug] = page
		context.Raw[slug] = data
		context.Slugs = append(context.Slugs, slug)
	}

	slice.Sort(context.Slugs, func(i, j int) bool {
		return context.Slugs[i] < context.Slugs[j]
	})

	context.Nav = mapping.EntryToIndexItem(index.Nav)
}
