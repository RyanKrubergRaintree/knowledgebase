// +build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var skipdirs = map[string]bool{
	".git": true,
	"cmd":  true,
}

func main() {
	filename := fmt.Sprintf("kb-%s.zip", time.Now().Format("2006-01-02-15-04"))
	log.Println(filename)
	file, err := os.Create(filepath.Join("..", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	AddDir(w, ".")
}

func Add(w *zip.Writer, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	hdr := &zip.FileHeader{}
	hdr.Name = path
	hdr.UncompressedSize = uint32(info.Size())
	hdr.SetMode(info.Mode())
	hdr.SetModTime(info.ModTime())

	body, err := w.CreateHeader(hdr)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(body, file); err != nil {
		log.Fatal(err)
	}
}

func AddDir(w *zip.Writer, dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && skipdirs[path] {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, ".exe") {
			return nil
		}

		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		Add(w, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
