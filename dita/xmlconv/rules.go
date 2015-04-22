package xmlconv

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type Encoder interface {
	EncodeToken(token interface{}) error
	Flush() error
}

type Callback func(enc Encoder, dec *xml.Decoder, start *xml.StartElement) error

type Customizer struct {
	Element   func(start *xml.StartElement) (skip bool, err error)
	CharData  func(char xml.CharData) (skip bool, err error)
	Comment   func(comment xml.Comment) (skip bool, err error)
	ProcInst  func(procinst *xml.ProcInst) (skip bool, err error)
	Directive func(directive xml.Directive) (skip bool, err error)
}

type Rules struct {
	Translate map[string]string
	Callback  map[string]Callback
	Unwrap    map[string]bool
	Remove    map[string]bool

	Handle Customizer
}

func NewRules() *Rules {
	return &Rules{
		Translate: map[string]string{},
		Callback:  map[string]Callback{},
		Unwrap:    map[string]bool{},
		Remove:    map[string]bool{},
		Handle:    Customizer{},
	}
}

func (rules *Rules) ConvertBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := rules.Convert(bytes.NewReader(data), &buf)
	return buf.Bytes(), err
}

func (rules *Rules) ConvertString(text string) (string, error) {
	var buf bytes.Buffer
	err := rules.Convert(bytes.NewBufferString(text), &buf)
	return buf.String(), err
}

func (rules *Rules) Convert(input io.Reader, output io.Writer) error {
	dec, enc := xml.NewDecoder(input), NewHTMLEncoder(output)

	dec.AutoClose = xml.HTMLAutoClose
	dec.Entity = xml.HTMLEntity

	err := rules.Encode(enc, dec, nil)
	if err != nil {
		return err
	}
	return enc.Flush()
}

func (rules *Rules) ConvertAny(enc Encoder, dec *xml.Decoder, token interface{}) error {
	switch token := token.(type) {
	case xml.StartElement:
		return rules.ConvertElement(enc, dec, &token)
	case xml.EndElement:
		return fmt.Errorf("unexpected end element at %d", dec.InputOffset())
	case xml.CharData:
		tokencopy := token.Copy()
		if rules.Handle.CharData != nil {
			if skip, err := rules.Handle.CharData(tokencopy); err != nil || skip {
				return err
			}
		}
		return enc.EncodeToken(tokencopy)
	case xml.Comment:
		tokencopy := token.Copy()
		if rules.Handle.Comment != nil {
			if skip, err := rules.Handle.Comment(tokencopy); err != nil || skip {
				return err
			}
		}
		return enc.EncodeToken(tokencopy)
	case xml.ProcInst:
		tokencopy := token.Copy()
		if rules.Handle.ProcInst != nil {
			if skip, err := rules.Handle.ProcInst(&tokencopy); err != nil || skip {
				return err
			}
		}
		return enc.EncodeToken(tokencopy)
	case xml.Directive:
		tokencopy := token.Copy()
		if rules.Handle.Directive != nil {
			if skip, err := rules.Handle.Directive(tokencopy); err != nil || skip {
				return err
			}
		}
		return enc.EncodeToken(tokencopy)
	default:
		panic("invalid token")
	}
}

func (rules *Rules) Encode(enc Encoder, dec *xml.Decoder, start *xml.StartElement) error {
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
		if err := rules.ConvertAny(enc, dec, token); err != nil {
			return err
		}
	}
}

func (rules *Rules) unwrap(enc Encoder, dec *xml.Decoder, start *xml.StartElement) error {
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
		if err := rules.ConvertAny(enc, dec, token); err != nil {
			return err
		}
	}
}

func (rules *Rules) ConvertElement(enc Encoder, dec *xml.Decoder, start *xml.StartElement) error {
	name := start.Name.Local
	if fn, ok := rules.Callback[name]; ok {
		return fn(enc, dec, start)
	}

	if remove := rules.Remove[name]; remove {
		return dec.Skip()
	}

	if unwrap := rules.Unwrap[name]; unwrap {
		return rules.unwrap(enc, dec, start)
	}

	startcopy := start.Copy()

	if rules.Handle.Element != nil {
		if skip, err := rules.Handle.Element(&startcopy); err != nil || skip {
			if skip {
				dec.Skip()
			}
			return err
		}
	}

	if newname, ok := rules.Translate[startcopy.Name.Local]; ok {
		startcopy.Name.Local = newname
	}

	if err := enc.EncodeToken(startcopy); err != nil {
		return err
	}

	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		if end, done := token.(xml.EndElement); done {
			endcopy := end
			if endcopy.Name != start.Name {
				return fmt.Errorf("invalid end token: start:%v end:%v", start, end)
			}
			if newname, ok := rules.Translate[endcopy.Name.Local]; ok {
				endcopy.Name.Local = newname
			}
			return enc.EncodeToken(endcopy)
		}

		if err := rules.ConvertAny(enc, dec, token); err != nil {
			return err
		}
	}
}

func (rules *Rules) Merge(with *Rules) {
	if with.Translate != nil {
		for name, result := range with.Translate {
			rules.Translate[name] = result
		}
	}
	if with.Callback != nil {
		for name, cb := range with.Callback {
			rules.Callback[name] = cb
		}
	}
	if with.Unwrap != nil {
		for name, unwrap := range with.Unwrap {
			rules.Unwrap[name] = unwrap
		}
	}
	if with.Remove != nil {
		for name, remove := range with.Remove {
			rules.Remove[name] = remove
		}
	}
	if with.Handle.Element != nil {
		rules.Handle.Element = with.Handle.Element
	}
	if with.Handle.CharData != nil {
		rules.Handle.CharData = with.Handle.CharData
	}
	if with.Handle.Comment != nil {
		rules.Handle.Comment = with.Handle.Comment
	}
	if with.Handle.ProcInst != nil {
		rules.Handle.ProcInst = with.Handle.ProcInst
	}
	if with.Handle.Directive != nil {
		rules.Handle.Directive = with.Handle.Directive
	}
}
