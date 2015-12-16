package dita

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

func EnsureAudience(data []byte) ([]byte, error) {
	var out bytes.Buffer
	enc := xml.NewEncoder(&out)
	dec := xml.NewDecoder(bytes.NewReader(data))
	err := ensureAudience(enc, dec, xml.StartElement{})
	enc.Flush()
	return out.Bytes(), err
}

func ensureAudience(enc *xml.Encoder, dec *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if end, done := token.(xml.EndElement); done {
			enc.EncodeToken(token)
			if end.Name != start.Name {
				return fmt.Errorf("invalid end token at %d: start:%v end:%v", dec.InputOffset(), start, end)
			}
			return nil
		}

		if start, ok := token.(xml.StartElement); ok {
			audience := getAttr(start, "audience")
			if audience == "html" || audience == "print" {
				if err := dec.Skip(); err != nil {
					return err
				}
				continue
			}
			if getAttr(start, "print") == "printonly" {
				if err := dec.Skip(); err != nil {
					return err
				}
				continue
			}

			enc.EncodeToken(token)
			if err := ensureAudience(enc, dec, start); err != nil {
				return err
			}
			continue
		}

		enc.EncodeToken(token)
	}
	return nil
}
