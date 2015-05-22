package pgdb

import (
	"database/sql"
	"fmt"

	"github.com/raintreeinc/knowledgebase/kbserver"
)

var _ kbserver.Database = &Database{}

type Database struct{ *sql.DB }

func New(params string) (*Database, error) {
	sdb, err := sql.Open("postgres", params)
	if err != nil {
		return nil, fmt.Errorf("failed to load DB: %s", err)
	}
	db := &Database{sdb}
	return db, db.Initialize()
}

func (db *Database) exec(query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	return err
}

func (db *Database) reset() error {
	return db.exec(`
		-- Reset DB state;
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		GRANT ALL ON SCHEMA public TO postgres;
		GRANT ALL ON SCHEMA public TO public;
		COMMENT ON SCHEMA public IS 'standard public schema';
	`)
}

func (db *Database) Initialize() error {
	_ = db.reset()
	return db.exec(`
		CREATE TABLE Groups (
			Name    TEXT   PRIMARY KEY UNIQUE,
			Public  BOOL   NOT NULL, -- can be read by everyone

			Description TEXT NOT NULL DEFAULT ''
		);

		CREATE TABLE Users (
			Name  TEXT   PRIMARY KEY UNIQUE,
			Email TEXT   NOT NULL,
			
			Description TEXT NOT NULL DEFAULT ''
		);

		CREATE TABLE Memberships (
			UserName  TEXT NOT NULL REFERENCES Users(Name),
			GroupName TEXT NOT NULL REFERENCES Groups(Name),

			CONSTRAINT Memberships_PKEY PRIMARY KEY (UserName, GroupName)			
		);

		CREATE TABLE Pages (
			Owner     TEXT  NOT NULL REFERENCES Groups(Name),
			Slug      TEXT  NOT NULL,
			Data      JSONB NOT NULL,
			Version   INT   NOT NULL DEFAULT 0,
			
			Tags      TEXT[] NOT NULL DEFAULT '{}',

			Created  TIMESTAMP NOT NULL DEFAULT current_timestamp,
			Modified TIMESTAMP NOT NULL DEFAULT current_timestamp,

			CONSTRAINT Pages_PKEY PRIMARY KEY (Owner, Slug)
		);
		CREATE INDEX PagesTagIndex ON Pages USING gin ((Data->'tags'));
		CREATE INDEX PagesSynopsisIndex ON Pages USING gin ((Data->'synopsis'));

		-- Triggers to automatically update modified date
		CREATE FUNCTION UpdateModifiedDate() RETURNS TRIGGER AS $$
		BEGIN
			NEW.Modified := NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE PLPGSQL;
		CREATE TRIGGER PagesUpdateModifiedDate
		BEFORE UPDATE ON Pages
		FOR EACH ROW EXECUTE PROCEDURE UpdateModifiedDate();
	`)
}
