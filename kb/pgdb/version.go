package pgdb

import (
	"database/sql"
	"fmt"
)

type migration struct {
	Version int
	Scripts []string
}

var migrations = []*migration{
	{
		Version: 1,
		Scripts: []string{
			`CREATE TYPE Rights
				AS ENUM ('blocked', 'reader', 'editor', 'moderator')`,
			`CREATE TABLE Groups (
				ID      TEXT PRIMARY KEY,
				OwnerID TEXT NOT NULL REFERENCES Groups(ID),
				Name    TEXT NOT NULL,
				Public  BOOL NOT NULL,

				Description TEXT NOT NULL DEFAULT ''
			)`,
			`CREATE TABLE Users (
				ID      TEXT PRIMARY KEY,
				Name    TEXT NOT NULL,
				Email   TEXT NOT NULL,
				Company TEXT NOT NULL DEFAULT '',
				Admin   BOOL NOT NULL DEFAULT false,

				AuthID       TEXT NOT NULL,
				AuthProvider TEXT NOT NULL,

				MaxAccess Rights NOT NULL DEFAULT 'moderator'
			)`,
			`CREATE TABLE Membership (
				GroupID  TEXT   NOT NULL REFERENCES Groups(ID),
				UserID   TEXT   NOT NULL REFERENCES Users(ID),

				CONSTRAINT Membership_PKEY PRIMARY KEY (GroupID, UserID)
			)`,
			`CREATE TABLE Community (
				GroupID   TEXT   NOT NULL REFERENCES Groups(ID),
				MemberID  TEXT   NOT NULL REFERENCES Groups(ID),
				Access    Rights NOT NULL DEFAULT 'blocked',

				CONSTRAINT Community_PKEY PRIMARY KEY (GroupID, MemberID)
			)`,
			`CREATE TABLE Pages (
				OwnerID TEXT NOT NULL REFERENCES Groups(ID),
				Slug    TEXT NOT NULL PRIMARY KEY,

				Data    JSONB NOT NULL,
				Version INT   NOT NULL DEFAULT 0,

				Tags      TEXT[] NOT NULL DEFAULT '{}',
				TagSlugs  TEXT[] NOT NULL DEFAULT '{}',

				Created  TIMESTAMP NOT NULL DEFAULT current_timestamp,
				Modified TIMESTAMP NOT NULL DEFAULT current_timestamp
			)`,
			`CREATE TABLE PageJournal (
				Actor    TEXT   NOT NULL,
				Slug     TEXT   NOT NULL,
				Version  INT    NOT NULL,
				Action   TEXT NOT NULL,
				Data     JSONB  NOT NULL,
				Date     TIMESTAMP NOT NULL DEFAULT current_timestamp
			)`,
		},
	},
	{
		Version: 2,
		Scripts: []string{
			`CREATE TABLE GuestLogin (
				AuthID TEXT NOT NULL,
				Name   TEXT NOT NULL,
				Email  TEXT NOT NULL,
				Salt   BYTEA NOT NULL,
				DK     BYTEA NOT NULL,

				CONSTRAINT GuestLogin_PKEY PRIMARY KEY (Name)
			)`,
		},
	},
	{
		Version: 3,
		Scripts: []string{
			// Add caching columns
			`ALTER TABLE Pages
				ADD COLUMN Content  TSVECTOR,
				ADD COLUMN Title    TEXT NOT NULL DEFAULT '',
				ADD COLUMN Synopsis TEXT NOT NULL DEFAULT ''`,
		},
	},
	{
		Version: 4,
		Scripts: []string{
			// add automatic lookup field update function
			`CREATE FUNCTION Pages_Update() RETURNS trigger AS
			$$
				DECLARE
					Story TEXT;
				BEGIN
					SELECT INTO Story string_agg(Item.Content, ' ')
						FROM (SELECT CAST(jsonb_array_elements(new.Data->'story')->'text' AS TEXT) AS Content) Item;
					new.Content :=
						setweight(to_tsvector('english', coalesce(CAST(new.Data->>'title' AS TEXT),'')), 'A') ||
						setweight(to_tsvector('english', coalesce(CAST(new.Data->>'synopsis' AS TEXT),'')), 'B') ||
						setweight(to_tsvector('english', story), 'C');
					new.Title    := coalesce(CAST(new.Data->>'title' AS TEXT), '');
					new.Synopsis := coalesce(CAST(new.Data->>'synopsis' AS TEXT), '');
					RETURN new;
				END
			$$ LANGUAGE plpgsql VOLATILE
			COST 100`,
			// add trigger for changing pages
			`CREATE TRIGGER Pages_UpdateTrigger
				BEFORE INSERT OR UPDATE
				ON Pages FOR EACH ROW EXECUTE PROCEDURE Pages_Update()`,
			// refresh all the content using the Pages_Update
			`UPDATE Pages SET OwnerID = OwnerID`,
			// create index for content
			`CREATE INDEX PagesContentGIN ON Pages USING gin(Content)`,
		},
	},
	{
		Version: 5,
		Scripts: []string{
			`CREATE VIEW AccessView AS
				WITH Accesses AS (
					-- public pages
					SELECT Groups.ID AS GroupID, Users.ID AS UserID, 'reader'::Rights AS Access
					FROM Groups
					CROSS JOIN Users
					WHERE Groups.Public = true
				UNION ALL
					-- member of group
					SELECT Membership.GroupID, Membership.UserID, 'moderator'::Rights AS Rights
					FROM Membership
				UNION ALL
					-- member of group owner
					SELECT Groups.ID, Membership.UserID, 'moderator'::Rights AS Rights
					FROM Groups
					JOIN Membership ON Membership.GroupID = Groups.OwnerID
				UNION ALL
					-- member of group community
					SELECT Groups.ID, Membership.UserID, Community.Access
					FROM Groups
					JOIN Community ON Community.GroupID = Groups.ID
					JOIN Membership ON Membership.GroupID = Community.MemberID
				)
			SELECT Accesses.GroupID, Accesses.UserID, LEAST(MAX(Accesses.Access), Users.MaxAccess) AS Access
			FROM Accesses
			JOIN Users ON Users.ID = Accesses.UserID
			GROUP BY Accesses.GroupID, Accesses.UserID, Users.ID
			ORDER BY Accesses.GroupID, Accesses.UserID;`,
		},
	},
}

func (db *Database) createVersionTable() error {
	_, err := db.Exec(`
	DO $$ BEGIN
		CREATE TABLE IF NOT EXISTS Versions (
			Version INT NOT NULL UNIQUE,
			Updated TIMESTAMP NOT NULL DEFAULT current_timestamp
		);

		IF NOT EXISTS (SELECT 1 FROM Versions WHERE VERSION = 0) then
			INSERT INTO Versions (Version) VALUES (0);
		END IF;
	END; $$;`)
	if err != nil {
		return fmt.Errorf("Failed to create Version table: %v", err)
	}
	return nil
}

func (db *Database) isOldVersioning() bool {
	exists := func(table string) bool {
		err := db.QueryRow("SELECT FROM " + table).Scan()
		if err == sql.ErrNoRows {
			return false
		}
		return err == nil
	}
	return !exists(`Versions`) && exists(`Pages`)
}

func (db *Database) migrateFromOld() error {
	if err := db.createVersionTable(); err != nil {
		return err
	}

	err := db.migrate(&migration{
		Version: 2,
		Scripts: []string{
			`DROP TRIGGER Pages_UpdateTrigger ON Pages`,
			`DROP FUNCTION Pages_Update()`,
			`DROP INDEX PagesContentGIN`,
			`ALTER TABLE Pages DROP COLUMN Content`,
		},
	})
	if err != nil {
		db.Exec("DROP TABLE Versions")
	}
	return err
}

func (db *Database) migrate(mig *migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, q := range mig.Scripts {
		if _, err := tx.Exec(q); err != nil {
			return fmt.Errorf("Migration to version %d failed at step %d: %v", mig.Version, i+1, err)
		}
	}

	_, err = tx.Exec(`INSERT INTO Versions (Version) VALUES ($1)`, mig.Version)
	if err != nil {
		return fmt.Errorf("Version update failed: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Migration to version %d failed: %v", mig.Version, err)
	}

	return nil
}

func (db *Database) Initialize() error {
	if db.isOldVersioning() {
		if err := db.migrateFromOld(); err != nil {
			return err
		}
	} else {
		if err := db.createVersionTable(); err != nil {
			return err
		}
	}

	version := 0
	if err := db.QueryRow(`SELECT MAX(Version) FROM Versions`).Scan(&version); err != nil {
		return err
	}

	for _, mig := range migrations {
		if mig.Version > version {
			err := db.migrate(mig)
			if err != nil {
				return err
			}
			version = mig.Version
		}
	}
	return nil
}
