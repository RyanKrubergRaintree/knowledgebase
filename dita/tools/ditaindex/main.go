// This is a tool for testing the building of ditamap index and links
package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/raintreeinc/knowledgebase/dita/ditaindex"
)

type TopicByName []*ditaindex.Topic

func (a TopicByName) Len() int           { return len(a) }
func (a TopicByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TopicByName) Less(i, j int) bool { return a[i].Filename < a[j].Filename }

func printcaption(name string) {
	fmt.Println()
	fmt.Println("==================")
	fmt.Printf("= %-14s =\n", name)
	fmt.Println("==================")
	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "USAGE:")
		fmt.Fprintln(os.Stderr, "  nav root.ditamap")
		os.Exit(1)
	}

	index, errs := ditaindex.Load(os.Args[1])
	if len(errs) > 0 {
		printcaption("ERRORS")
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	printcaption("NAV")
	printnav(index.Nav, "", "    ")

	printcaption("TOPICS")

	topics := make([]*ditaindex.Topic, 0, len(index.Topics))
	for _, topic := range index.Topics {
		topics = append(topics, topic)
	}

	sort.Sort(TopicByName(topics))

	for _, topic := range topics {
		fmt.Println()
		fmt.Println("====", topic.Title, "==== ", topic.Filename)
		for _, links := range topic.Links {
			if links.Parent != nil {
				fmt.Print("\t^  ", links.Parent.Title, "\n")
			}
			if links.Prev != nil || links.Next != nil {
				fmt.Print("\t")
				if links.Prev != nil {
					fmt.Print("<- ", links.Prev.Title)
				}
				if links.Next != nil {
					if links.Prev != nil {
						fmt.Print(" - ")
					}
					fmt.Print(links.Next.Title, " ->")
				}
				fmt.Println()
			}
			if len(links.Children) > 0 {
				fmt.Print("\tv  ")
				for _, child := range links.Children {
					fmt.Print(child.Title, " ")
				}
				fmt.Println()
			}
			if len(links.Siblings) > 0 {
				fmt.Print("\t~  ")
				for _, sibling := range links.Siblings {
					fmt.Print(sibling.Title, " ")
				}
				fmt.Println()
			}
		}
	}
}

func printnav(e *ditaindex.Entry, prefix, indent string) {
	link := ""
	if e.Topic != nil {
		link = ">"
	}
	if len(e.Children) == 0 {
		fmt.Printf("%s- %s %s\n", prefix, e.Title, link)
		return
	}

	fmt.Printf("%s+ %s %s\n", prefix, e.Title, link)
	for _, child := range e.Children {
		printnav(child, prefix+indent, indent)
	}
}
