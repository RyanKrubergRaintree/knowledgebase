-- Reset DB state;
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
COMMENT ON SCHEMA public IS 'standard public schema';

-- Initialize
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
	GroupName TEXT NOT NULL REFERENCES Groups(Name)
);

CREATE TABLE Pages (
	Slug      TEXT  PRIMARY KEY,
	GroupName TEXT   NOT NULL REFERENCES Groups(Name),
	Data      JSONB NOT NULL,
	Version   INT   NOT NULL DEFAULT 0,

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
INSERT INTO Groups VALUES ('Community', True, '');
INSERT INTO Groups VALUES ('Engineering', False, 'Raintree Engineering');
INSERT INTO Groups VALUES ('Help', True, 'Raintree Help');

INSERT INTO Users VALUES ('Egon', 'egonelbre@gmail.com', 'Raintree Systems Inc.');

INSERT INTO Memberships Values ('Egon', 'Community');
INSERT INTO Memberships Values ('Egon', 'Engineering');

INSERT INTO Pages VALUES (
	'/home',
	'Community',
	'{"title":"world", "tags":["alpha", "beta"], "synopsis":"Lorem ipsum dolor sit amet, consectetur adipisicing elit."}',
	0
);

INSERT INTO Pages VALUES (
	'/work',
	'Engineering',
	'{"title":"hello", "tags":["beta", "gamma"], "synopsis":"Excepteur sint occaecat cupidatat non proident."}',
	0
);

-- -- Filter based on what user can modify
-- WHERE GroupName IN (SELECT GroupName FROM Memberships WHERE User = 'Egon')
-- -- And what is public
--    OR GroupName IN (SELECT Name FROM Groups WHERE Public = True)

-- select based on tag
SELECT
	Slug, GroupName, Data->'synopsis'
FROM Pages;

-- select based on tag
SELECT
	Slug, GroupName, Data, Version
FROM Pages
WHERE data @> '{"tags": ["gamma"]}';

-- get all tags
EXPLAIN SELECT
	jsonb_array_elements_text(data->'tags') as Tag,
	count(*) as Count
FROM Pages
WHERE GroupName IN (SELECT GroupName FROM Memberships WHERE User = 'Egon')
   OR GroupName IN (SELECT Name FROM Groups WHERE Public = True)
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