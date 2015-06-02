// This is a tool for testing the naming of dita page conversion
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/raintreeinc/knowledgebase/ditaconv"
	"github.com/raintreeinc/knowledgebase/kb"
)

var (
	outputdir = flag.String("dir", "out", "output directory for json files")
	verbose   = flag.Bool("v", false, "verbose output")

	maxconvert = flag.Int("max", -1, "maximum `number` of topics to convert")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "USAGE:")
		fmt.Fprintln(os.Stderr, "  ditaconvert root.ditamap")
		os.Exit(1)
	}

	index, errs := ditaconv.LoadIndex(args[0])
	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, "Loading errors:")
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	mapping, errs := ditaconv.CreateMapping(index)
	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, "Mapping errors:")
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	os.MkdirAll(*outputdir, 0755)

	topics := mapping.TopicsSorted()
	for i, topic := range topics {
		if (*maxconvert > 0) && i >= *maxconvert {
			break
		}

		slug := mapping.ByTopic[topic]

		if *verbose {
			fmt.Printf("%04d/%04d  %-40s %-40s\n", i+1, len(topics), slug, topic.Filename)
		}

		page, fatal, errs := mapping.Convert(topic)
		if fatal != nil {
			fmt.Fprintln(os.Stderr, "!!!", topic.Filename, fatal)
		}

		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintln(os.Stderr, "!", topic.Filename, err)
			}
		}

		if fatal == nil {
			data, err := json.MarshalIndent(page, "", "\t")
			if err != nil {
				fmt.Fprintln(os.Stderr, "!!", topic.Filename, err)
			}

			outfile := filepath.Join(*outputdir, slugToFilename(slug))
			err = ioutil.WriteFile(outfile, data, 0755)
			if err != nil {
				fmt.Fprintln(os.Stderr, "!!", topic.Filename, err)
			}
		}
	}
}

func slugToFilename(slug kb.Slug) string {
	return strings.Replace(string(slug), "/", "-", -1)
}
