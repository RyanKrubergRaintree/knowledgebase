package admin

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/raintreeinc/kbdita/ditaconv"
	"github.com/raintreeinc/kbdita/ditaindex"

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

	out, err := ioutil.TempDir("", "dita-")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Unzipping new dita:", out)
	defer os.Remove(out)

	unzip(out, bytes.NewReader(full.Bytes()), fileSize)

	index, errs := ditaindex.Load(filepath.Join(out, "dita", "contents.ditamap"))
	if len(errs) > 0 {
		fmt.Fprintf(w, "Loading: %v\n", errs)
	}

	mapping, errs := ditaconv.CreateMapping(index)
	if len(errs) > 0 {
		fmt.Fprintf(w, "Mapping: %v\n", errs)
	}

	c := session.DB("").C("Help")

	for slug, topic := range mapping.BySlug {
		page, fatal, errs := mapping.Convert(topic)
		page.Slug = "/" + page.Slug
		if fatal != nil {
			fmt.Fprintf(w, "Fatal: %v -> %v\n", slug, fatal)
			continue
		} else if len(errs) > 0 {
			fmt.Fprintf(w, "Convert: %v -> %v\n", slug, errs)
		}

		if err := c.Insert(page); err != nil {
			fmt.Fprintf(w, "Insert: %v -> %v\n", slug, errs)
		}
	}

	fmt.Fprintf(w, "\n\nCompleted.")
}

func unzip(dest string, src io.ReaderAt, size int64) error {
	r, err := zip.NewReader(src, size)
	if err != nil {
		return err
	}
	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
