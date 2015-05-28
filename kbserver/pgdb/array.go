// stringSlice support based on
//   https://gist.github.com/adharris/4163702
//   and comments
// Remove when Array support lands in pq
package pgdb

import (
	"bytes"
	"database/sql/driver"
	"encoding/csv"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type stringSlice []string

// for replacing escaped quotes except if it is preceded by a literal backslash
//  eg "\\" should translate to a quoted element whose value is \

var quoteEscapeRegex = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

// Scan convert to a slice of strings
// http://www.postgresql.org/docs/9.1/static/arrays.html#ARRAYS-IO
func (s *stringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	str := string(asBytes)

	// change quote escapes for csv parser
	str = quoteEscapeRegex.ReplaceAllString(str, `$1""`)
	str = strings.Replace(str, `\\`, `\`, -1)
	// remove braces
	str = str[1 : len(str)-1]
	csvReader := csv.NewReader(strings.NewReader(str))

	slice, err := csvReader.Read()

	if err != nil {
		return err
	}

	(*s) = stringSlice(slice)

	return nil
}

func (s stringSlice) Value() (driver.Value, error) {
	if s == nil {
		return []byte{}, nil
	}

	var buffer bytes.Buffer

	buffer.WriteString("{")
	last := len(s) - 1
	for i, val := range s {
		buffer.WriteString(strconv.Quote(val))
		if i != last {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}
