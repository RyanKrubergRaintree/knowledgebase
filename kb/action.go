package kb

import (
	"fmt"
	"time"
)

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
	m, ismap := (item).(map[string]interface{})
	if !ismap {
		return nil, false
	}
	return (Item)(m), true
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