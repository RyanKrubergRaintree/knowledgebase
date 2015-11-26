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

func splitref(ref string) (file string, path []string) {
	if !strings.Contains(ref, "#") {
		return ref, nil
	}

	tokens := strings.SplitN(ref, "#", 2)
	return tokens[0], strings.Split(tokens[1], "/")
}

func sameroot(a, b []string) bool {
	if len(a) == 0 || len(a) != len(b) {
		return false
	}

	for i, v := range a[:len(a)-1] {
		if !strings.EqualFold(v, b[i]) {
			return false
		}
	}
	return true
}

func (r *replacer) conref(start xml.StartElement) error {
	startfile, startpath := splitref(getAttr(start, "conref"))
	endfile, endpath := splitref(getAttr(start, "conrefend"))

	if endfile == "" && len(endpath) == 0 {
		endfile, endpath = startfile, startpath
	}

	if startfile == "" && endfile == "" {
		startfile = r.file
		endfile = r.file
	} else {
		startfile = filepath.Join(filepath.Dir(r.file), filepath.FromSlash(startfile))
		endfile = filepath.Join(filepath.Dir(r.file), filepath.FromSlash(endfile))
	}

	if startfile != endfile {
		return errors.New("conref and conrefend are in different files")
	}
	if !sameroot(startpath, endpath) {
		return errors.New("conref and conrefend have different root elements")
	}
	if len(startpath) == 0 || len(endpath) == 0 {
		return errors.New("invalid conref path")
	}

	data, err := os.Open(startfile)
	if err != nil {
		return fmt.Errorf("Problem opening %v: %v", startfile, err)
	}

	defer data.Close()
	dec := xml.NewDecoder(data)

	s, err := walknodepath(dec, startpath)
	if err != nil {
		if err == io.EOF {
			return errors.New("did not find conref")
		}
		return err
	}

	lastid := endpath[len(endpath)-1]
	for {
		r.enc.EncodeToken(s)
		err := (&replacer{startfile, r.enc}).emit(dec, s)
		r.enc.EncodeToken(xml.EndElement{Name: s.Name})
		if err != nil {
			return err
		}

		if strings.EqualFold(lastid, getAttr(s, "id")) {
			return nil
		}

		for {
			x, err := dec.Token()
			if err != nil {
				if err == io.EOF {
					return errors.New("did not find conrefend")
				}
				return err
			}

			var ok bool
			s, ok = x.(xml.StartElement)
			if ok {
				break
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
