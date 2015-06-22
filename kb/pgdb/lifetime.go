package pgdb

import "fmt"

var setup = []string{`
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'rights') THEN
			CREATE TYPE Rights AS ENUM ('blocked', 'reader', 'editor', 'moderator');
		END IF;
	END$$;
`, `
	CREATE TABLE IF NOT EXISTS
	Groups (
		ID      TEXT PRIMARY KEY,
		OwnerID TEXT NOT NULL REFERENCES Groups(ID),
		Name    TEXT NOT NULL,
		Public  BOOL NOT NULL,

		Description TEXT NOT NULL DEFAULT ''
	)
`, `
	CREATE TABLE IF NOT EXISTS
	Users (
		ID      TEXT PRIMARY KEY,
		Name    TEXT NOT NULL,
		Email   TEXT NOT NULL,
		Company TEXT NOT NULL DEFAULT '',
		Admin   BOOL NOT NULL DEFAULT false,

		AuthID       TEXT NOT NULL,
		AuthProvider TEXT NOT NULL
	)
`, `
	-- describes user members
	CREATE TABLE IF NOT EXISTS
	Membership (
		GroupID  TEXT   NOT NULL REFERENCES Groups(ID),
		UserID   TEXT   NOT NULL REFERENCES Users(ID),

		CONSTRAINT Membership_PKEY PRIMARY KEY (GroupID, UserID)
	)
`, `
	-- describes community groups
	CREATE TABLE IF NOT EXISTS
	Community (
		GroupID   TEXT   NOT NULL REFERENCES Groups(ID),
		MemberID  TEXT   NOT NULL REFERENCES Groups(ID),
		Access    Rights NOT NULL DEFAULT 'blocked',

		CONSTRAINT Community_PKEY PRIMARY KEY (GroupID, MemberID)
	)
`, `
	CREATE TABLE IF NOT EXISTS
	Pages (
		OwnerID TEXT NOT NULL REFERENCES Groups(ID),
		Slug    TEXT NOT NULL PRIMARY KEY,

		Data    JSONB NOT NULL,
		Version INT   NOT NULL DEFAULT 0,

		Tags      TEXT[] NOT NULL DEFAULT '{}',
		TagSlugs  TEXT[] NOT NULL DEFAULT '{}',

		Created  TIMESTAMP NOT NULL DEFAULT current_timestamp,
		Modified TIMESTAMP NOT NULL DEFAULT current_timestamp
	)
`, `
	CREATE TABLE IF NOT EXISTS
	PageJournal (
		Actor    TEXT   NOT NULL,
		Slug     TEXT   NOT NULL,
		Version  INT    NOT NULL,
		Action   TEXT NOT NULL,
		Data     JSONB  NOT NULL,
		Date     TIMESTAMP NOT NULL DEFAULT current_timestamp
	)
`, `
-- this function updates Pages.Content Vector
CREATE OR REPLACE FUNCTION Pages_Update() RETURNS trigger AS
$$
	DECLARE
		Story TEXT;
	BEGIN
		SELECT INTO Story string_agg(Item.Content, ' ')
			FROM (SELECT CAST(jsonb_array_elements(new.Data->'story')->'text' AS TEXT) AS Content) Item;
		new.Content :=
			setweight(to_tsvector('english', coalesce(CAST(new.Data->'title' AS TEXT),'')), 'A') ||
			setweight(to_tsvector('english', coalesce(CAST(new.Data->'synopsis' AS TEXT),'')), 'B') ||
			setweight(to_tsvector('english', story), 'C');
		RETURN new;
	END
$$ LANGUAGE plpgsql VOLATILE
COST 100;
`, `
DO $$
BEGIN
	BEGIN
		ALTER TABLE Pages ADD COLUMN Content TSVECTOR;

		CREATE TRIGGER Pages_UpdateTrigger
		BEFORE INSERT OR UPDATE
		ON Pages FOR EACH ROW EXECUTE PROCEDURE Pages_Update();

		-- do a nop instruction to add all items
		UPDATE Pages SET OwnerID = OwnerID;

		CREATE INDEX PagesContentGIN ON Pages USING gin(Content);
	EXCEPTION
		WHEN duplicate_column THEN RAISE NOTICE 'Content column already exists.';
	END;
END;
$$
`, `
	CREATE TABLE IF NOT EXISTS
	GuestLogin (
		AuthID TEXT NOT NULL,
		Name   TEXT NOT NULL,
		Email  TEXT NOT NULL,
		Salt   BYTEA NOT NULL,
		DK     BYTEA NOT NULL,

		CONSTRAINT GuestLogin_PKEY PRIMARY KEY (Name)
	)
`, `
DO $$
BEGIN
	BEGIN
		ALTER TABLE Users ADD COLUMN MaxAccess Rights NOT NULL DEFAULT 'moderator';
	EXCEPTION
		WHEN duplicate_column THEN RAISE NOTICE 'MaxAccess column already exists.';
	END;
END;
$$
`}

func (db *Database) Initialize() error {
	for i, q := range setup {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("%d: %v\n-QUERY-\n%v\n--\n", i, err, q)
		}
	}
	return nil
}
