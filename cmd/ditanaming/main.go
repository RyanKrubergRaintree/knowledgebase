// This is a tool for testing the naming of dita pages with slugs
package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/raintreeinc/knowledgebase/dita/ditaconv"
	"github.com/raintreeinc/knowledgebase/dita/ditaindex"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "USAGE:")
		fmt.Fprintln(os.Stderr, "  ditanaming root.ditamap")
		os.Exit(1)
	}

	index, errs := ditaindex.Load(os.Args[1])
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

	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 0, '\t', 0)
	for slug, topic := range mapping.BySlug {
		fmt.Fprintln(tw, slug, "\t", topic.Title, "\t", topic.Filename)
	}
	tw.Flush()
}
