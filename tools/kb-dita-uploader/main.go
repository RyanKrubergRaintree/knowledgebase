package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/raintreeinc/knowledgebase/ditaconv"
	"github.com/raintreeinc/knowledgebase/extra/ditaindex"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kb/pgdb"
	"github.com/raintreeinc/knowledgebase/module/dita"
)

var (
	configfile = flag.String("config", "kb-dita-uploader.json", "configuration file")
	stoponerr  = flag.Bool("stop", false, "don't upload if there are problems in converting")
	killonerr  = flag.Bool("kill", false, "don't try upload other mappings")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	config := &Config{}

	file, err := os.Create(time.Now().Format("upload-2006-01-02T150405.log"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(io.MultiWriter(file, os.Stdout))

	if err := config.ReadFromFile(*configfile); err != nil {
		log.Println(err)
		return
	}

	onlyupload := flag.Args()
	if len(onlyupload) == 0 {
		for _, m := range config.Mapping {
			onlyupload = append(onlyupload, m.Group)
		}
	}

	allStart := time.Now()
	for _, name := range onlyupload {
		log.Println()
		log.Println()
		log.Println("========================================")
		log.Println("====", name)
		start := time.Now()
		if err := Upload(name, config); err != nil {
			log.Println("ERROR:", err)
			if *killonerr {
				return
			}
		}
		log.Println("==== Completed in ", time.Since(start))
		log.Println("========================================")
	}

	log.Println("==== Everything completed in ", time.Since(allStart))
}

func fileexists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func Upload(name string, config *Config) error {
	var p *CopyParams
	for _, cp := range config.Mapping {
		if strings.EqualFold(name, cp.Group) {
			p = &cp
			break
		}
	}

	if p == nil {
		return fmt.Errorf("no mapping named ", name)
	}

	log.Println()
	log.Println("== Loading index:", p.Ditamap)
	if !fileexists(p.Ditamap) {
		return errors.New("ditamap doesn't exist")
	}

	index, errs := ditaconv.LoadIndex(p.Ditamap)
	if len(errs) > 0 {
		log.Println()
		log.Println(errs)
		if *stoponerr {
			return errors.New("errors in index")
		}
	}

	log.Println()
	log.Println("== Creating mapping")
	mapping, errs := ditaconv.CreateMapping(index)
	if len(errs) > 0 {
		log.Println()
		log.Println(errs)
		if *stoponerr {
			return errors.New("errors in mapping")
		}
	}

	owner := kb.Slugify(p.Group)
	for topic, slug := range mapping.ByTopic {
		ownerslug := owner + "=" + slug
		mapping.ByTopic[topic] = ownerslug
		delete(mapping.BySlug, slug)
		mapping.BySlug[ownerslug] = topic
	}

	navindex := ditaindex.EntryToItem(mapping, index.Nav)

	mapping.Rules.Merge(dita.RaintreeDITA())

	log.Println()
	log.Println("== Converting pages")
	log.Println()
	pages := make(map[kb.Slug]*kb.Page)
	for _, topic := range mapping.BySlug {
		page, fatal, errs := mapping.Convert(topic)
		if fatal != nil {
			log.Println(fatal)
			if *stoponerr {
				return errors.New("error in converting " + topic.Title)
			}
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

	navindexslug := kb.Slug(owner + "=index")
	pages[navindexslug] = &kb.Page{
		Slug:     navindexslug,
		Title:    "Index",
		Synopsis: "Help navigation index",
		Story: kb.Story{
			ditaindex.New("index", navindex),
		},
	}

	log.Println()
	log.Println("== Connecting to DB")
	DB, err := pgdb.New(config.ConnectionParams())
	if err != nil {
		return err
	}
	defer DB.Close()

	// Try create group if it does not exist
	err = DB.Context("admin").Groups().Create(kb.Group{
		ID:          owner,
		OwnerID:     owner,
		Name:        p.Group,
		Public:      true,
		Description: p.Description,
	})
	if err == nil {
		log.Println("== Created group:", owner)
	} else if err != kb.ErrGroupExists {
		log.Println("== Creating group: ", owner)
		return err
	}

	log.Println()
	log.Println("== Uploading")
	log.Println()

	complete := 0
	total := len(pages)
	err = DB.Context("admin").Pages(owner).BatchReplace(pages, func(slug kb.Slug) {
		complete++
		log.Printf("%04d/%04d : %v\n", complete, total, slug)
	})

	return err
}

type CopyParams struct {
	Group       string
	Ditamap     string
	Description string
}

type Config struct {
	// all db params at once
	DBParams string

	// RDS is clearer setup for Amazon when DBParams is not defined
	RDS struct {
		User   string
		Pass   string
		DBName string
		Host   string
		Port   string
	}

	// Mapping contains how to map dita files to kb pages and groups
	// additionally it creates all the necessary groups if they don't already
	// exist
	Mapping []CopyParams
}

func (c *Config) LoadEnv() {
	c.DBParams = os.Getenv("DATABASE")
	c.RDS.User = os.Getenv("RDS_USERNAME")
	c.RDS.Pass = os.Getenv("RDS_PASSWORD")
	c.RDS.DBName = os.Getenv("RDS_DB_NAME")
	c.RDS.Host = os.Getenv("RDS_HOSTNAME")
	c.RDS.Port = os.Getenv("RDS_PORT")
}

func (c *Config) ReadFrom(r io.Reader) error {
	return json.NewDecoder(r).Decode(c)
}

func (c *Config) ReadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return c.ReadFrom(file)
}

func (c *Config) ConnectionParams() string {
	if c.DBParams != "" {
		return c.DBParams
	}

	return fmt.Sprintf(
		"user='%s' password='%s' dbname='%s' host='%s' port='%s'",
		c.RDS.User, c.RDS.Pass, c.RDS.DBName, c.RDS.Host, c.RDS.Port,
	)
}
