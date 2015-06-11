package cmds

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/raintreeinc/knowledgebase/ditaconv"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/module/dita"
)

func init() {
	Register(Command{
		Name: "upgrade-help",
		Desc: "For batch updating help content.",
		Run:  HelpUpgrade,
	})
}

func HelpUpgrade(DB kb.Database, fs *flag.FlagSet, args []string) {
	ditamap := fs.String("ditamap", "", "root ditamap")
	fs.Parse(args)

	if os.Getenv("DITAMAP") != "" {
		*ditamap = os.Getenv("DITAMAP")
	}

	if *ditamap == "" {
		flag.Usage()
	}

	index, errs := ditaconv.LoadIndex(*ditamap)
	log.Println(errs)

	mapping, errs := ditaconv.CreateMapping(index)
	log.Println(errs)

	owner := kb.Slug("help")
	for topic, slug := range mapping.ByTopic {
		ownerslug := owner + ":" + slug
		mapping.ByTopic[topic] = ownerslug
		delete(mapping.BySlug, slug)
		mapping.BySlug[ownerslug] = topic
	}

	mapping.Rules.Merge(dita.RaintreeDITA())

	pages := make(map[kb.Slug]*kb.Page)
	for _, topic := range mapping.BySlug {
		page, fatal, errs := mapping.Convert(topic)
		if fatal != nil {
			log.Println(fatal)
			continue
		} else if len(errs) > 0 {
			log.Println(errs)
		}

		if page.Slug == "" {
			log.Printf("No slug for \"%s\".", page.Title)
			continue
		}

		if page.Slug[0] == '/' {
			page.Slug = page.Slug[1:]
		}
		pages[page.Slug] = page
	}

	complete := 0
	total := len(pages)

	err := DB.Context("admin").Pages("help").BatchReplace(pages, func(slug kb.Slug) {
		complete++
		fmt.Printf("%04d/%04d : %v\n", complete, total, slug)
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Completed.")
	}
}
