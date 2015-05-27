# Knowledge Base Client

The Client is structured as:

- Convert: contains utilities for converting between URL, Link, Slug

Knowledge Base parts:

- KB.Page: handles page content and operations on the content
- KB.Stage: handles uploading, downloading pages. Staging area for page modifications.
- KB.Lineup: contains the list of stages currently open.
- KB.Crumbs: updates the location-bar hashes

View of KB:

- View.App: main layout of the page
- View.Stage: shows a single KB.Stage with meta information and loading
- View.Lineup: displays stages side-by-side
