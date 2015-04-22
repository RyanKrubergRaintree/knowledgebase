package xmlconv

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestZeroRules(t *testing.T) {
	cases := []string{
		"<xml></xml>",
		"<p><p><p><p></p><br/></p></p></p>",
		"<p>Hello, World</p>",
		"<p>Hello, <b>World</b></p>",
		"<p>Hello, <b>World<br/></b></p>",
		"<p>Hello,<!-- Comment --><b>World<br/></b></p>",
	}

	rules := NewRules()
	for _, input := range cases {
		output, err := rules.ConvertString(input)
		if err != nil {
			t.Error(err)
		}
		expected := strings.Replace(input, "<br/>", "<br>", -1)
		if output != expected {
			t.Errorf("invalid output: got %v expected %v", output, expected)
		}
	}
}

func TestRemove(t *testing.T) {
	in := `<main><remove><p><remove/></p><x></x></remove></main>`
	expected := `<main></main>`

	rules := NewRules()
	rules.Remove["remove"] = true

	output, err := rules.ConvertString(in)
	if err != nil {
		t.Error(err)
	}
	if output != expected {
		t.Errorf("invalid output: got %v expected %v", output, expected)
	}
}

func TestUnwrap(t *testing.T) {
	in := `<main><unwrap><p>X<unwrap/>Y</p>Y<x></x></unwrap></main>`
	expected := `<main><p>XY</p>Y<x></x></main>`

	rules := NewRules()
	rules.Unwrap["unwrap"] = true

	output, err := rules.ConvertString(in)
	if err != nil {
		t.Error(err)
	}
	if output != expected {
		t.Errorf("invalid output: got %v expected %v", output, expected)
	}
}

func TestCallback(t *testing.T) {
	in := `<main><cb x="Y"><a /><b /></cb></main>`
	expected := `<main><q w="Z"></q></main>`

	rules := NewRules()
	rules.Callback["cb"] = func(enc Encoder, dec *xml.Decoder, start *xml.StartElement) error {
		err := enc.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: "q"},
			Attr: []xml.Attr{{Name: xml.Name{Local: "w"}, Value: "Z"}},
		})
		if err != nil {
			return err
		}
		err = enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "q"}})
		if err != nil {
			return err
		}
		return dec.Skip()
	}

	rules.Unwrap["unwrap"] = true

	output, err := rules.ConvertString(in)
	if err != nil {
		t.Error(err)
	}
	if output != expected {
		t.Errorf("invalid output: got %v expected %v", output, expected)
	}
}
