package kbdita

import (
	"net/http"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

var _ *kb.Page
var _ kbserver.System = &System{}

type System struct {
	name    string
	ditamap string
	server  *kbserver.Server
}

func New(name, ditamap string, server *kbserver.Server) *System {
	sys := &System{
		name:    name,
		ditamap: ditamap,
		server:  server,
	}
	sys.init()
	return sys
}

func (sys *System) Name() string { return sys.name }

func (sys *System) init() {
	//
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
