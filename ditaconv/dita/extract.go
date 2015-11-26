package dita

import (
	"encoding/xml"
	"errors"
	"strings"
)

func walknodepath(dec *xml.Decoder, unmatched []string) (xml.StartElement, error) {
	if len(unmatched) == 0 {
		return xml.StartElement{}, errors.New("invalid nodepath")
	}

	for {
		token, err := dec.Token()
		if err != nil {
			return xml.StartElement{}, err
		}

		if start, ok := token.(xml.StartElement); ok {
			id := getAttr(start, "id")
			if strings.EqualFold(id, unmatched[0]) {
				unmatched = unmatched[1:]
			}
			if len(unmatched) == 0 {
				return start, nil
			}
		}
	}

	panic("unreachable")
}
