package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kb/pgdb"
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

	allStart := time.Now()
	if err := removehelp(config); err != nil {
		log.Println(err)
	}
	if err := fixSlug(config); err != nil {
		log.Println(err)
	}
	log.Println("==== Everything completed in ", time.Since(allStart))
}

func fixSlug(config *Config) error {
	DB, err := pgdb.New(config.ConnectionParams())
	if err != nil {
		return err
	}
	defer DB.Close()

	rows, err := DB.Query(`SELECT Slug, OwnerID, Title, Data FROM Pages`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var slug, ownerid, title string
		var data []byte

		err := rows.Scan(&slug, &ownerid, &title, &data)
		if err != nil {
			return err
		}

		log.Println("Processing:", slug)
		newslug := kb.Slugify(ownerid + "=" + title)

		page := &kb.Page{}
		if err := json.Unmarshal(data, page); err != nil {
			return err
		}

		page.Slug = newslug
		page.Synopsis = kb.ExtractSynopsis(page)
		tags := kb.ExtractTags(page)
		tagSlugs := kb.SlugifyTags(tags)

		newdata, err := json.Marshal(page)
		if err != nil {
			return fmt.Errorf("failed to serialize page: %v", err)
		}

		_, err = DB.Exec(`
			UPDATE Pages
			SET Slug = $2,
				Data = $3,
				Tags = $4,
				TagSlugs = $5,
				Version = Version
			WHERE Slug = $1
		`, slug, string(newslug), newdata, stringSlice(tags), stringSlice(tagSlugs))

		if err != nil {
			return err
		}
	}

	return nil
}

func removehelp(config *Config) error {
	DB, err := pgdb.New(config.ConnectionParams())
	if err != nil {
		return err
	}
	defer DB.Close()

	_, err = DB.Exec(`DELETE FROM Pages WHERE OwnerID LIKE 'help-%'`)
	return err
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
