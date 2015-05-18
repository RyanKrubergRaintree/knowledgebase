package kb

import "time"

type PageEntry struct {
	Owner    string    `json:"owner"`
	Slug     Slug      `json:"slug"`
	Title    string    `json:"title"`
	Synopsis string    `json:"synopsis"`
	Tags     []string  `json:"tags"`
	Modified time.Time `json:"modified"`
}

type TagEntry struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
