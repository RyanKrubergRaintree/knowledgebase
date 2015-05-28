package pgdb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/raintreeinc/knowledgebase/kb"
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

func (db *Database) Destroy() error {
	return db.exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		GRANT ALL ON SCHEMA public TO postgres;
		GRANT ALL ON SCHEMA public TO public;
		COMMENT ON SCHEMA public IS 'standard public schema';
	`)
}

func (db *Database) Initialize() error {
	//TODO: fix this setup
	err := db.exec(`
		CREATE TABLE Groups (
			ID      TEXT   PRIMARY KEY,
			Name    TEXT   NOT NULL,
			Public  BOOL   NOT NULL,

			Description TEXT NOT NULL DEFAULT ''
		);

		CREATE TABLE Users (
			ID    TEXT   PRIMARY KEY,
			Name  TEXT   NOT NULL,
			Email TEXT   NOT NULL,
			
			Description TEXT NOT NULL DEFAULT ''
		);

		CREATE TABLE Memberships (
			UserID  TEXT NOT NULL REFERENCES Users(ID),
			GroupID TEXT NOT NULL REFERENCES Groups(ID),

			CONSTRAINT Memberships_PKEY PRIMARY KEY (UserID, GroupID)			
		);

		CREATE TABLE Pages (
			Owner     TEXT  NOT NULL REFERENCES Groups(ID), -- e.g. community
			Slug      TEXT  NOT NULL PRIMARY KEY, -- e.g. community:welcome-visitor
			Data      JSONB NOT NULL,
			Version   INT   NOT NULL DEFAULT 0,
			
			Tags      TEXT[] NOT NULL DEFAULT '{}',
			NormTags  TEXT[] NOT NULL DEFAULT '{}',

			Created  TIMESTAMP NOT NULL DEFAULT current_timestamp,
			Modified TIMESTAMP NOT NULL DEFAULT current_timestamp
		);

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
	log.Println(err)
	return nil
}

func (db *Database) CanWrite(user, group kb.Slug) (result bool) {
	err := db.QueryRow(`SELECT
		EXISTS(SELECT FROM Memberships WHERE UserID = $1 AND GroupID = $2)
	`, user, group).Scan(&result)
	return (err == nil) && result
}

func (db *Database) CanRead(user, group kb.Slug) (result bool) {
	err := db.QueryRow(`SELECT
		(SELECT Public FROM Groups WHERE ID = $2)
     	OR EXISTS( SELECT FROM Memberships WHERE UserID = $1 AND GroupID = $2 )
	`, user, group).Scan(&result)
	return (err == nil) && result
}
