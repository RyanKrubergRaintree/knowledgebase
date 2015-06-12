select to_tsvector('english', 
	coalesce(cast(Data->'title' AS TEXT), '')
) from pages;

select Data->'story'
from pages
limit 50;

select Data#>>'{story,1,text}'
from pages
limit 50;

select json_each('{"story": [{"text":"alpha"}, {"text":"beta"}, {"text":"gamma"}]}'::jsonb->'story')

SELECT Slug 
FROM Pages
WHERE to_tsvector('english', coalesce(cast(Data->'title' AS TEXT),'')) @@ plainto_tsquery('english', 'testing')
   OR to_tsvector('english', coalesce(cast(Data->'story' AS TEXT),'')) @@ plainto_tsquery('english', 'testing')
LIMIT 100;

SELECT Slug 
FROM PagesSearch
WHERE Title_TS @@ plainto_tsquery('english', 'testing')
   OR Story_TS @@ plainto_tsquery('english', 'testing')
LIMIT 100;


CREATE INDEX PagesGIN ON Pages USING gin (Data jsonb_path_ops);
CREATE INDEX PagesFullTextGIN ON PAGES USING GIN(
	to_tsvector('english',
		coalesce(cast(Data->'title' AS TEXT),'') || ' ' ||
		coalesce(cast(Data->'story' AS TEXT), ''))
);
DROP INDEX PagesFullTextGIn;

CREATE INDEX PagesTitleGIN ON PAGES USING GIN(to_tsvector('english',coalesce(cast(Data->'title' AS TEXT),'')));
CREATE INDEX PagesStoryGIN ON PAGES USING GIN(to_tsvector('english',coalesce(cast(Data->'story' AS TEXT),'')));
DROP INDEX PagesTitleGIN;
DROP INDEX PagesStoryGIN;


CREATE VIEW PagesSearch AS
	SELECT 
		OwnerID, Slug,
		coalesce(cast(Data->'title' AS TEXT),'') AS Title, 
		coalesce(cast(Data->'synopsis' AS TEXT),'') AS Synopsis, 
		to_tsvector('english',coalesce(cast(Data->'title' AS TEXT),'')) AS Title_TS,
		to_tsvector('english',coalesce(cast(Data->'story' AS TEXT),'')) AS Story_TS
	FROM Pages;

SELECT * FROM PagesSearch LIMIT 100;

select to_tsvector('english', 
	coalesce(cast(Data->'story' AS TEXT), '')
) from pages;