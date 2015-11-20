package pgdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"

	_ "github.com/lib/pq"
)

var ErrNotImplemented = errors.New("not implemented")

var _ kb.Database = &Database{}

type Database struct {
	*sql.DB
}

func New(params string) (*Database, error) {
	sdb, err := sql.Open("postgres", params)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err)
	}

	db := &Database{DB: sdb}
	return db, nil
}

func (db Access) BoolQuery(q string, args ...interface{}) bool {
	err := db.QueryRow(q, args...).Scan()
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Println(err)
	}
	return err == nil
}

func (db Database) Context(user kb.Slug) kb.Context { return Context{db, user} }

type Context struct {
	Database
	ActiveUser kb.Slug
}

func (ctx Context) ActiveUserID() kb.Slug { return ctx.ActiveUser }
func (ctx Context) Access() kb.Access     { return Access{ctx} }
func (ctx Context) Users() kb.Users       { return Users{ctx} }
func (ctx Context) Groups() kb.Groups     { return Groups{ctx} }

func (ctx Context) GuestLogin() kb.GuestLogin { return GuestLogin{ctx} }

func (ctx Context) Index(user kb.Slug) kb.Index  { return Index{ctx, user} }
func (ctx Context) Pages(group kb.Slug) kb.Pages { return Pages{ctx, group} }

func dupkey(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}

func (ctx Context) pageEntries(filter string, args ...interface{}) (entries []kb.PageEntry, err error) {
	rows, err := ctx.Query(`
	SELECT
		Slug,
		Title,
		Synopsis,
		Tags,
		Modified
	FROM Pages
	`+filter, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry kb.PageEntry

		xtags := stringSlice{}
		err := rows.Scan(
			&entry.Slug,
			&entry.Title,
			&entry.Synopsis,
			&xtags,
			&entry.Modified,
		)
		entry.Tags = []string(xtags)

		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}
