package ditaconv

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/raintreeinc/knowledgebase/ditaconv/dita"
)

type Dir string

// converts path to filepath
func (d Dir) Filepath(name string) string {
	return filepath.Join(string(d), filepath.FromSlash(name))
}

func (d Dir) Open(name string) (file io.ReadCloser, err error) {
	file, err = os.Open(d.Filepath(name))
	return
}

func (d Dir) ReadFile(name string) (data []byte, err error) {
	return ioutil.ReadFile(d.Filepath(name))
}

func (d Dir) Stat(name string) (os.FileInfo, error) {
	return os.Stat(d.Filepath(name))
}

func (d Dir) LoadTopic(filename string) (*Topic, error) {
	topic, err := dita.LoadTopic(d.Filepath(filename))
	if err != nil {
		return &Topic{
			Filename: filename,
			Title:    trimext(filepath.Base(filename)),
		}, fmt.Errorf("reading \"%s\": %s", filename, err)
	}

	return &Topic{
		Filename:   filename,
		Title:      first(topic.NavTitle, topic.Title),
		ShortTitle: topic.Title,
		Synopsis:   topic.ShortDesc,

		Original: topic,
	}, nil
}

func (d Dir) LoadMap(filename string) (*Map, error) {
	data, err := d.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("loading \"%s\": %s", filename, err)
	}

	node := &dita.MapNode{}
	if err := xml.Unmarshal(data, &node); err != nil {
		return nil, fmt.Errorf("reading \"%s\": %s", filename, err)
	}

	return &Map{Node: node}, nil
}

func first(xs ...string) string {
	for _, x := range xs {
		if x != "" {
			return x
		}
	}
	return ""
}

func trimext(name string) string {
	return name[0 : len(name)-len(filepath.Ext(name))]
}

func canonicalName(name string) string {
	if runtime.GOOS == "windows" {
		return strings.ToLower(name)
	}
	return name
}
