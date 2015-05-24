-- Reset DB state;
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
COMMENT ON SCHEMA public IS 'standard public schema';

-- -- Filter based on what user can modify
-- WHERE GroupName IN (SELECT GroupName FROM Memberships WHERE User = 'Egon')
-- -- And what is public
--    OR GroupName IN (SELECT Name FROM Groups WHERE Public = True)

-- simple select
SELECT
	Slug, Owner, Data->'synopsis'
FROM Pages;

-- select based on tag
SELECT
	Owner, Slug,
	Data->'title' as Title,
	Data->'synopsis' as Synopsis,
	Tags,
	Modified
FROM Pages
WHERE Tags @> ARRAY['gamma'];

SELECT
	Name, Email, Description, array_agg(Memberships.GroupName) as Groups
FROM Users
JOIN Memberships ON (Users.Name = Memberships.UserName)
GROUP BY Users.Name;


SELECT
	Owner, Slug,
	Data->'title' as Title,
	Data->'synopsis' as Synopsis,
	Tags,
	Modified
FROM Pages
WHERE (Tags @> ARRAY['beta']) 
  AND (    Owner IN (SELECT Name      FROM Groups      WHERE Public = TRUE)
	OR Owner IN (SELECT GroupName FROM Memberships WHERE User = 'Egon'))
ORDER BY Owner, Slug;

SELECT (SELECT Public FROM Groups WHERE Name = 'xxx')
     OR EXISTS(SELECT FROM Memberships WHERE UserName = 'Egon' AND GroupName = 'xxx')
	

EXISTS(
	
) OR (SELECT)

-- get all tags
EXPLAIN
SELECT
	unnest(Tags) as Tag,
	count(*) as Count
FROM Pages
WHERE Owner IN (SELECT GroupName FROM Memberships WHERE User = 'Egon')
   OR Owner IN (SELECT Name FROM Groups WHERE Public = True)
GROUP BY Tag
;

SELECT current_user;


SELECT * FROM Pages
WHERE
EXISTS (
	SELECT FROM Memberships WHERE UserName = 'Egon' AND GroupName = 'engineerin'
);


SELECT * FROM Pages
IF EXISTS (SELECT 1 FROM Memberships WHERE UserName = 'Egon' AND GroupName = 'engineerin')
;

UPDATE Pages SET Data = '{}'
WHERE Slug = '/work'
  AND EXISTS (SELECT FROM Memberships WHERE UserName = 'Egon' AND GroupName = 'engineerin');

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