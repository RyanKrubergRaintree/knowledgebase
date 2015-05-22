package kb

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"unicode"
)

// Slug is a string where Slugify(string(slug)) == slug
type Slug string

func (slug *Slug) Scan(value interface{}) error {
	*slug = Slug(value.(string))
	return nil
}

func (slug Slug) Value() (driver.Value, error) {
	return string(slug), nil
}

// ValidateSlug verifies whether a `slug` is valid
func ValidateSlug(slug Slug) error {
	if len(slug) == 0 {
		return fmt.Errorf("slug cannot be empty")
	}

	conv := Slugify(string(slug))
	if slug != conv {
		return fmt.Errorf(`slugification modified the slug`)
	}

	return nil
}

// Slugify converts text to a slug
//
// * numbers, '/' are left intact
// * letters will be lowercased (if possible)
// * '-', ',', '.', ' ', '_' will be converted to '-'
// * other symbols or punctuations will be converted to html entity reference name
//   (if there exists such reference name)
// * everything else will be converted to '-'
//
// Example:
//   "&Hello_世界/+!" ==> "amp-hello-世界/plus-excl"
//   "Hello  World  /  Test" ==> "hello-world/test"
func Slugify(s string) Slug {
	cutdash := true
	emitdash := false

	slug := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsNumber(r) || unicode.IsLetter(r) {
			if emitdash && !cutdash {
				slug = append(slug, '-')
			}
			slug = append(slug, unicode.ToLower(r))

			emitdash = false
			cutdash = false
			continue
		}
		switch r {
		case '/':
			slug = append(slug, r)
			emitdash = false
			cutdash = true
		case '-', ',', '.', ' ', '_':
			emitdash = true
		default:
			if name, exists := runename[r]; exists {
				if !cutdash {
					slug = append(slug, '-')
				}
				slug = append(slug, []rune(name)...)
				cutdash = false
			}
			emitdash = true
		}
	}

	if len(slug) == 0 {
		return Slug("-")
	}

	return Slug(slug)
}

func SplitOwner(url string) (owner string, page Slug) {
	start := 0
	if strings.HasPrefix(url, "/") {
		start = 1
	}

	i := strings.Index(url, ":")
	if i < 0 {
		return "", Slugify(url[start:])
	}
	return url[start:i], Slugify(url[i+1:])
}
