// +build ignore

package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"

	"github.com/raintreeinc/knowledgebase/ditaconv/dita"
)

func main() {
	flag.Parse()

	topic, err := dita.LoadTopic(flag.Arg(0))
	if err != nil {
		log.Println(err)
		return
	}

	enc := xml.NewEncoder(os.Stdout)
	defer enc.Flush()

	enc.Encode(topic)
}
