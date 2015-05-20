package kbserver

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jaschaephraim/lrserver"
)

//TODO: extract this into a separate library

type Source struct {
	IsDev bool
	Dir   string

	LiveReload *lrserver.Server
}

func NewSource(dir string, isdev bool) *Source {
	source := &Source{
		IsDev: isdev,
		Dir:   dir,
	}
	if source.IsDev {
		var err error
		source.LiveReload, err = lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			// Delay live reload start to avoid port collision
			time.Sleep(1 * time.Second)
			log.Fatal(source.LiveReload.ListenAndServe())
		}()

		go source.monitorChanges()
	}
	return source
}

func (s *Source) monitorChanges() {
	lastTime := time.Now()
	for range time.NewTicker(5 * time.Second).C {
		filepath.Walk(s.Dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".js") {
				return nil
			}

			if info.ModTime().After(lastTime) {
				path = strings.TrimPrefix(path, s.Dir)
				path = filepath.ToSlash(path)
				s.LiveReload.Reload(path)
			}
			return nil
		})
		lastTime = time.Now()
	}
}

func (s *Source) extractDeps(path string) []string {
	filename := filepath.Join(s.Dir, filepath.FromSlash(path))

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Extracting deps from %s: %s", path, err)
		return nil
	}

	deps := []string{}
	rx := regexp.MustCompile(`//\s*import\s+"([^"]*)"`)
	for _, match := range rx.FindAllStringSubmatch(string(file), -1) {
		deps = append(deps, match[1])
	}

	return deps
}

func (s *Source) sorted(sources []string) ([]string, error) {
	sort.Strings(sources)

	result := []string{}
	deps := make(map[string][]string)
	for _, file := range sources {
		deps[file] = s.extractDeps(file)
	}

	for file, deps := range deps {
		log.Println("\t", file, "->", deps)
	}
	sorted := make(map[string]bool)

	// brute force topological sort
	for pass := 0; pass < 100; pass += 1 {
		changes := false

	nextpkg:
		for _, file := range sources {
			if sorted[file] {
				continue
			}

			for _, dep := range deps[file] {
				if !sorted[dep] {
					continue nextpkg
				}
			}

			changes = true
			sorted[file] = true
			result = append(result, file)
		}

		if len(sorted) == len(sources) {
			return result, nil
		}

		if !changes {
			return nil, errors.New("sources contain a circular dependency")
		}
	}

	return nil, errors.New("sources contain a circular dependency")
}

func (s *Source) Files() []string {
	r := []string{}

	err := filepath.Walk(s.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".js") {
			return nil
		}

		path = strings.TrimPrefix(path, s.Dir)
		r = append(r, filepath.ToSlash(path))
		return nil
	})

	if err != nil {
		log.Println(err)
		return nil
	}

	files, err := s.sorted(r)
	if err != nil {
		log.Println(err)
		return nil
	}

	return files
}

func (s *Source) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, SafeFile(s.Dir, r.URL.Path))
}
