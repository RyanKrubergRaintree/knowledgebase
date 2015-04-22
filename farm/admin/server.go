package admin

import (
	"fmt"
	"net/http"
)

type Server struct {
	Database string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, AdminPage)
	case "/updatehelp":
		s.updateHelp(w, r)
	default:
		http.NotFound(w, r)
	}
}

const AdminPage = `
<html>
<title>Admin Page</title>
<body>
	<form action="/updatehelp" method="post" enctype="multipart/form-data">
		<label for="file">Filename:</label>
		<input type="file" name="file" id="file">
		<input type="submit" name="submit" value="Submit">
	</form>
</body>
</html>
`
