package admin

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/egonelbre/fedwiki"

	"gopkg.in/mgo.v2-unstable"
	"gopkg.in/mgo.v2-unstable/bson"
)

func (s *Server) updateHelp(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(s.Database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	if true { // drop old entries
		session.DB("").C("Help").RemoveAll(bson.M{})
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var full bytes.Buffer
	fileSize, err := full.ReadFrom(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file.Close()

	z, err := zip.NewReader(bytes.NewReader(full.Bytes()), fileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := session.DB("").C("Help")
	uploadPage := func(f *zip.File) error {
		if f.FileInfo().IsDir() {
			return nil
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		var page fedwiki.Page
		dec := json.NewDecoder(rc)
		if err := dec.Decode(&page); err != nil {
			return err
		}

		if err := c.Insert(page); err != nil {
			return err
		}

		return nil
	}
	log.Println(len(z.File))

	for _, file := range z.File {
		err := uploadPage(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "Completed.")
}
