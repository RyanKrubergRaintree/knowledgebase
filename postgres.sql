-- Reset DB state;
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
COMMENT ON SCHEMA public IS 'standard public schema';

-- Initialize
CREATE TABLE Groups (
	ID      SERIAL PRIMARY KEY,
	Name    TEXT   NOT NULL UNIQUE,
	Public  BOOL   NOT NULL, -- can be read by everyone

	Description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE Users (
	ID    SERIAL PRIMARY KEY,
	Name  TEXT   NOT NULL UNIQUE,
	Email TEXT   NOT NULL,
	
	Description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE Memberships (
	UserID  INT NOT NULL REFERENCES Users(ID),
	GroupID INT NOT NULL REFERENCES Groups(ID)
);

CREATE TABLE Pages (
	Slug     TEXT  PRIMARY KEY,
	GroupID  INT   NOT NULL REFERENCES Groups(ID),
	Data     JSONB NOT NULL,
	Version  INT   NOT NULL DEFAULT 0,

	Created  TIMESTAMP NOT NULL DEFAULT current_timestamp,
	Modified TIMESTAMP NOT NULL DEFAULT current_timestamp
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

-- create some stub data
INSERT INTO Groups VALUES (0, 'Community', True, '');
INSERT INTO Groups VALUES (1, 'Engineering', False, 'Raintree Engineering');
INSERT INTO Groups VALUES (2, 'Help', True, 'Raintree Help');

INSERT INTO Users VALUES (0, 'Egon', 'egonelbre@gmail.com', 'Raintree Systems Inc.');

INSERT INTO Memberships Values (0, 0);
INSERT INTO Memberships Values (0, 1);

INSERT INTO Pages VALUES (
	'/home',
	0,
	'{"title":"world", "tags":["alpha", "beta"], "synopsis":"Lorem ipsum dolor sit amet, consectetur adipisicing elit."}',
	0
);

INSERT INTO Pages VALUES (
	'/work',
	1,
	'{"title":"hello", "tags":["beta", "gamma"], "synopsis":"Excepteur sint occaecat cupidatat non proident."}',
	0
);

-- -- Filter based on what user can modify
-- WHERE GroupID IN (SELECT GroupID FROM Memberships WHERE UserID = 0)
-- -- And what is public
--   OR GroupID IN (SELECT ID FROM Groups WHERE Public = True)

-- select based on tag
SELECT
	Slug, GroupID, Data->'synopsis'
FROM Pages;

-- select based on tag
SELECT
	Slug, GroupID, Data, Version
FROM Pages
WHERE data @> '{"tags": ["gamma"]}';

-- get all tags
EXPLAIN SELECT
	jsonb_array_elements_text(data->'tags') as Tag,
	count(*) as Count
FROM Pages
WHERE GroupID IN (SELECT GroupID FROM Memberships WHERE UserID = 0)
   OR GroupID IN (SELECT ID FROM Groups WHERE Public = True)
GROUP BY Tag;

-- all pages
SELECT * FROM Pages;

UPDATE Pages
SET Data = '{
	"tags":["beta", "gamma", "delta"],
	"title":"hello",
	"synopsis":"Hello world.",
	"items":[
		{"type":"text", "text": "Lorem"},
		{"type":"text", "text": "Lorem Ipsum"},
		{"type":"image", "data": "xyz"}
	]}'
WHERE Slug = '/work';

SELECT * FROM Pages;

-- full text search
--    http://www.postgresql.org/docs/9.4/static/textsearch-controls.html
--    http://www.postgresql.org/docs/9.4/static/functions-textsearch.html

-- example full text search
SELECT
	Slug
FROM Pages
WHERE	to_tsvector('english', 
		coalesce(cast(Data->'title' AS TEXT),'') || ' ' || 
		coalesce(cast(Data->'tags' AS TEXT),'') || ' ' || 
		coalesce(cast(Data->'items' AS TEXT), '')
	) @@ to_tsquery('english', 'lorem | alpha');

-- inspect the query
SELECT
	Slug,
	to_tsvector('english', 
		coalesce(cast(Data->'title' AS TEXT),'') || ' ' || 
		coalesce(cast(Data->'tags' AS TEXT),'') || ' ' || 
		coalesce(cast(Data->'items' AS TEXT), '')
	)
FROM Pages;

SELECT plainto_tsquery('english', 'The Fat Rats');