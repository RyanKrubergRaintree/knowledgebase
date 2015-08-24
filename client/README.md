# Knowledge Base Client

The Client is structured as:

- kb.convert: contains utilities for converting between URL, Link, Slug
- kb.Slugify: converts text to a possible slug page

Knowledge Base parts:

- kb.Page: handles page content and operations on the content
- kb.Stage: handles uploading, downloading pages. Staging area for page modifications.
- kb.Lineup: contains the list of stages currently open.
- kb.Crumbs: updates the location-bar hashes

View of KB:

- App: main layout of the page
- kb.Stage.View: shows a single KB.Stage with meta information and loading
- kb.Lineup.View: displays stages side-by-side
