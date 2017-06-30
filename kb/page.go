package kb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"
)

var (
	ErrUnknownAction = errors.New("unknown action")
)

// Page represents a federated wiki page
type Page struct {
	Version  int       `json:"version"`
	Slug     Slug      `json:"slug"`
	Title    string    `json:"title"`
	Synopsis string    `json:"synopsis,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
	Story    Story     `json:"story,omitempty"`
}

func (p *Page) Write(w io.Writer) error {
	data, err := json.Marshal(p)
	if err != nil {
		panic("unable to marshal page " + err.Error())
	}
	n, err := w.Write(data)
	if err == nil && n < len(data) {
		return io.ErrShortWrite
	}
	return err
}

// Story is the viewable content of the page
type Story []Item

// Apply modifies the page with an action
func (page *Page) Apply(action Action) error {
	fn, ok := actionfns[action.Type()]
	if !ok {
		return ErrUnknownAction
	}

	err := fn(page, action)
	if err != nil {
		return err
	}

	page.Version++
	return nil
}

// IndexOf returns the index of an item with `id`
// ok = false, if that item doesn't exist
func (s Story) IndexOf(id string) (index int, ok bool) {
	for i, item := range s {
		if item.ID() == id {
			return i, true
		}
	}
	return -1, false
}

// insertAt adds an `item` after position `i`
func (s *Story) insertAt(i int, item Item) {
	t := *s
	t = append(t, Item{})
	copy(t[i+1:], t[i:])
	t[i] = item
	*s = t
}

// Prepend adds the `item` as the first item in story
func (s *Story) Prepend(item Item) {
	s.insertAt(0, item)
}

// Appends adds the `item` as the last item in story
func (s *Story) Append(item ...Item) {
	*s = append(*s, item...)
}

// InsertAfter adds the `item` after the item with `id`
func (s *Story) InsertAfter(id string, item Item) error {
	if i, ok := s.IndexOf(id); ok {
		s.insertAt(i+1, item)
		return nil
	}
	return fmt.Errorf("invalid item id '%v'", id)
}

// SetByID replaces item with `id` with `item`
func (s Story) SetByID(id string, item Item) error {
	if i, ok := s.IndexOf(id); ok {
		s[i] = item
		return nil
	}
	return fmt.Errorf("invalid item id '%v'", id)
}

// Move moves the item with `id` after the item with `afterId`
func (ps *Story) Move(id string, afterId string) error {
	item, err := ps.RemoveByID(id)
	if err != nil {
		return err
	}
	if afterId != "" {
		return ps.InsertAfter(afterId, item)
	}
	ps.Prepend(item)
	return nil
}

// Removes item with `id`
func (s *Story) RemoveByID(id string) (item Item, err error) {
	if i, ok := s.IndexOf(id); ok {
		t := *s
		item = t[i]
		copy(t[i:], t[i+1:])
		t = t[:len(t)-1]
		*s = t
		return item, nil
	}
	return item, fmt.Errorf("missing item id '%v'", id)
}

func NewID() string { return fmt.Sprintf("%016x", rand.Int63()) }

// Item represents a federated wiki Story item
type Item map[string]interface{}

// Val returns a string value from key
func (item Item) Val(key string) string {
	if s, ok := item[key].(string); ok {
		return s
	}
	return ""
}

// Type returns the item `type`
func (item Item) Type() string { return item.Val("type") }

// ID returns the `item` identificator
func (item Item) ID() string { return item.Val("id") }

func ReadJSONPage(r io.Reader) (*Page, error) {
	dec := json.NewDecoder(r)
	page := &Page{}
	err := dec.Decode(page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func ReadJSONAction(r io.Reader) (Action, error) {
	dec := json.NewDecoder(r)
	action := make(Action)
	err := dec.Decode(&action)
	if err != nil {
		return nil, err
	}

	return action, nil
}

// Action represents a operation that can be applied to a fedwiki.Page
type Action map[string]interface{}

// Str returns string value by the key
// if that key doesn't exist, it will return an empty string
func (action Action) Str(key string) string {
	if s, ok := action[key].(string); ok {
		return s
	}
	return ""
}

// Type returns the action type attribute
func (action Action) Type() string {
	return action.Str("type")
}

// Item returns the item attribute
func (action Action) Item() (Item, bool) {
	item, ok := action["item"]
	if !ok {
		return nil, false
	}
	m, isitem := (item).(Item)
	if !isitem {
		m, ismap := (item).(map[string]interface{})
		if !ismap {
			return nil, false
		}
		return (Item)(m), true
	}
	return m, true
}

// Time returns the time when the action occurred
func (action Action) Time() (t time.Time, err error) {
	val, ok := action["time"]
	if !ok {
		return time.Time{}, fmt.Errorf("time not found")
	}
	switch val := val.(type) {
	case string:
		return time.Parse(time.RFC3339, val)
	case int: // assume unix timestamp
		return time.Unix(int64(val), 0), nil
	case int64: // assume unix timestamp
		return time.Unix(val, 0), nil
	}

	return time.Time{}, fmt.Errorf("unknown date format")
}

// actionfns defines how each action type is applied
var actionfns = map[string]func(p *Page, a Action) error{
	"add": func(p *Page, action Action) error {
		item, ok := action.Item()
		if !ok {
			return fmt.Errorf("no item in action")
		}

		after := action.Str("after")
		if after == "" {
			p.Story.Prepend(item)
			return nil
		}
		return p.Story.InsertAfter(after, item)
	},
	"edit": func(p *Page, action Action) error {
		item, ok := action.Item()
		if !ok {
			return fmt.Errorf("no item in action")
		}
		return p.Story.SetByID(action.Str("id"), item)
	},
	"remove": func(p *Page, action Action) error {
		_, err := p.Story.RemoveByID(action.Str("id"))
		return err
	},
	"move": func(p *Page, action Action) error {
		return p.Story.Move(action.Str("id"), action.Str("after"))
	},
	"create": func(p *Page, action Action) error {
		return nil
	},
	"fork": func(p *Page, action Action) error {
		return nil
	},
}
