package kb

import (
	"fmt"
	"time"
)

// Page represents a federated wiki page
type Page struct {
	Owner    string  `json:"owner"`
	Slug     Slug    `json:"slug"`
	Title    string  `json:"title"`
	Synopsis string  `json:"synopsis,omitempty"`
	Story    Story   `json:"story,omitempty"`
	Journal  Journal `json:"journal,omitempty"`
}

// Story is the viewable content of the page
type Story []Item

// Journal contains the history of the Page
type Journal []Action

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
	return nil
}

// LastModified returns the date when the page was last modified
// if there is no such date it will return a zero time
func (page *Page) LastModified() time.Time {
	for i := len(page.Journal) - 1; i >= 0; i-- {
		if t, err := page.Journal[i].Time(); err == nil && !t.IsZero() {
			return t
		}
	}
	return time.Time{}
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
func (s *Story) Append(item Item) {
	*s = append(*s, item)
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

// Item represents a federated wiki Story item
type Item map[string]interface{}

// Val returns a string value from key
func (item Item) Val(key string) string {
	if s, ok := item[key].(string); ok {
		return s
	}
	return ""
}

// ID returns the `item` identificator
func (item Item) ID() string { return item.Val("id") }
