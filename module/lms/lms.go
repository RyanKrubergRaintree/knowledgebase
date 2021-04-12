package lms

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/raintreeinc/knowledgebase/kb"
)

var _ kb.Module = &Module{}

type Module struct {
	server *kb.Server
	router *mux.Router
}

// New LMS module that acts as a limited LRS
func New(server *kb.Server) *Module {
	mod := &Module{
		server: server,
		router: mux.NewRouter(),
	}
	mod.init()
	return mod
}

// Info
func (mod *Module) Info() kb.Group {
	return kb.Group{
		ID:          "lms",
		Name:        "LMS",
		Public:      true,
		Description: "Learning managament system",
	}
}

func (mod *Module) init() {
	// create temp folder for uploads
	path, _ := os.Getwd()
	_ = os.Mkdir(filepath.FromSlash(path+"/temp/"), 666)
	mod.createUser()

	mod.router.HandleFunc("/lms=/uploadContent/", mod.listUploadedContent).Methods("GET")
	mod.router.HandleFunc("/lms=/uploadContent/", mod.uploadContent).Methods("POST")
}

//  Create default user for LMS uploads
func (mod *Module) createUser() {
	name := "lmsuser"
	_, err := mod.server.Database.Context("admin").Users().ByID(kb.Slugify(name))

	if err == kb.ErrUserNotExist {
		user := kb.User{
			AuthID:       name,
			AuthProvider: "guest",
			ID:           kb.Slugify(name),
			Email:        "lmsuser@raintreeinc.com",
			Name:         name,
			MaxAccess:    kb.Reader,
		}

		_ = mod.server.Database.Context("admin").Users().Create(user)
	}
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

// Pages
func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{{
		Slug:     "lms=lms",
		Title:    "LMS module.",
		Synopsis: "LMS module.",
	}}
}

func (mod *Module) listUploadedContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ListFilesFromBucket(w)
}

func (mod *Module) uploadContent(w http.ResponseWriter, r *http.Request) {
	err, fileNameWithPath := saveFileFromHttpRequestToServer(r)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	if uploadError, uploadedFilePath := uploadFileFromServerToS3(fileNameWithPath); uploadError == nil {
		fmt.Fprintf(w, uploadedFilePath)
	} else {
		fmt.Fprintf(w, uploadError.Error())
	}

	_ = os.Remove(fileNameWithPath)
}
