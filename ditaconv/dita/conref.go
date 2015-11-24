package dita

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getAttr(n xml.StartElement, key string) (val string) {
	for _, attr := range n.Attr {
		if attr.Name.Local == key {
			return attr.Value
		}
	}
	return ""
}

type replacer struct {
	file string
	enc  *xml.Encoder
}

func (r *replacer) conref(start xml.StartElement) error {
	conref := getAttr(start, "conref")
	if !strings.Contains(conref, "#") {
		return errors.New("Invalid conref value.")
	}

	tokens := strings.Split(conref, "#")
	file, nodepath := tokens[0], tokens[1]

	full := filepath.Join(filepath.Dir(r.file), filepath.FromSlash(file))
	if file == "" {
		full = r.file
	}

	data, err := os.Open(full)
	if err != nil {
		return fmt.Errorf("Problem opening %v: %v", full, err)
	}
	defer data.Close()
	dec := xml.NewDecoder(data)

	unmatched := strings.Split(nodepath, "/")
	if len(unmatched) == 0 {
		return errors.New("Invalid target node path.")
	}

	for {
		token, err := dec.Token()
		if err == io.EOF {
			return errors.New("Did not find conref " + conref)
		}
		if err != nil {
			return err
		}

		if start, ok := token.(xml.StartElement); ok {
			id := getAttr(start, "id")
			if !strings.EqualFold(id, unmatched[0]) {
				continue
			}
			unmatched = unmatched[1:]
			if len(unmatched) == 0 {
				r.enc.EncodeToken(start)
				err := (&replacer{full, r.enc}).emit(dec, start)
				r.enc.EncodeToken(xml.EndElement{Name: start.Name})
				return err
			}
		}
	}
}

func (r *replacer) emit(dec *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if end, done := token.(xml.EndElement); done {
			if end.Name != start.Name {
				return fmt.Errorf("invalid end token at %d: start:%v end:%v", dec.InputOffset(), start, end)
			}
			return nil
		}

		if start, ok := token.(xml.StartElement); ok {
			var err error
			if getAttr(start, "conref") != "" {
				err = r.conref(start)
				dec.Skip()
			} else {
				r.enc.EncodeToken(token)
				err = r.emit(dec, start)
				r.enc.EncodeToken(xml.EndElement{Name: start.Name})
			}
			if err != nil {
				return err
			}
			continue
		}

		r.enc.EncodeToken(token)
	}
	return nil
}

func replaceConrefs(data []byte, filename string) ([]byte, error) {
	if !bytes.Contains(data, []byte("conref=")) {
		return data, nil
	}

	var out bytes.Buffer
	enc := xml.NewEncoder(&out)

	dec := xml.NewDecoder(bytes.NewReader(data))
	err := (&replacer{filename, enc}).emit(dec, xml.StartElement{})

	enc.Flush()
	return out.Bytes(), err
}
