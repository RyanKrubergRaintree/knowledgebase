// +build ignore

package main

import (
	"archive/tar"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	output = flag.String("output", filepath.Join("..", "deploy.tar"), "output tar file")
)

func main() {
	file, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//gz := gzip.NewWriter(file)
	//defer gz.Close()

	t := tar.NewWriter(file)
	defer t.Close()

	AddDir(t, ".")
	os.Chdir("..")
	AddDir(t, "kbclient")
}

func Add(t *tar.Writer, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	hdr := &tar.Header{
		Name:    path,
		Size:    info.Size(),
		Mode:    int64(info.Mode()),
		ModTime: info.ModTime(),
	}
	if err := t.WriteHeader(hdr); err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(t, file); err != nil {
		log.Fatal(err)
	}
}

func AddDir(t *tar.Writer, dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".git") {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		Add(t, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
