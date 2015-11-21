package dita

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/raintreeinc/knowledgebase/ditaconv/xmlconv"
)

func ExtractTitle(filename, nodepath string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	dec := xml.NewDecoder(file)

	unmatched := strings.Split(nodepath, "/")
	if len(unmatched) == 0 {
		return "", errors.New("Invalid target node path.")
	}

	for {
		token, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return "", nil
			}
			return "", err
		}

		if start, ok := token.(xml.StartElement); ok {
			id := getAttr(start, "id")
			if !strings.EqualFold(id, unmatched[0]) {
				continue
			}
			unmatched = unmatched[1:]
			if len(unmatched) == 0 {
				return extractTitleTag(dec)
			}
		}
	}
}

func extractTitleTag(dec *xml.Decoder) (string, error) {
	for {
		token, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return "", nil
			}
			return "", err
		}
		if _, done := token.(xml.EndElement); done {
			return "", nil
		}

		if start, ok := token.(xml.StartElement); ok {
			if start.Name.Local == "title" {
				text, _ := xmlconv.Text(dec, nil)
				return text, nil
			}
			dec.Skip()
		}
	}
	return "", nil
}
