// +build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var ZIP *zip.Writer

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func run(name string, args ...string) error {
	fmt.Println("> ", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	cmd.Env = append([]string{
		"GOOS=linux",
		"GOARCH=amd64",
		"CGO_ENABLED=0",
	}, os.Environ()...)

	return cmd.Run()
}

func build() {
	check(run("go", "build", "-v", "-o", filepath.Join(".bin", "run"), "."))
	AddDir(".bin")
	AddDir(".ebextensions")

	AddGlob("*.json")
	AddGlob("Docker*")
	AddGlob("LICENSE-*")

	AddDir("client")
}

func main() {
	os.Mkdir(".bin", 0755)
	os.Mkdir(".deploy", 0755)

	filename := fmt.Sprintf("%s.zip", time.Now().Format("2006-01-02-15-04"))

	file, err := os.Create(filepath.Join(".deploy", filename))
	check(err)
	defer file.Close()

	fmt.Println("Creating:", filename)

	ZIP = zip.NewWriter(file)
	build()
	ZIP.Close()
}

// filename with forward slashes
func AddFile(filename string) {
	fmt.Printf("  %-40s", filename)
	defer fmt.Println("+")

	file, err := os.Open(filepath.FromSlash(filename))
	check(err)
	defer file.Close()

	w, err := ZIP.Create(filename)
	check(err)
	_, err = io.Copy(w, file)
	check(err)
}

// glob with forward slashes
func AddGlob(glob string) {
	fmt.Printf("G %v\n", glob)
	matches, err := filepath.Glob(filepath.FromSlash(glob))
	check(err)
	for _, match := range matches {
		AddFile(filepath.ToSlash(match))
	}
}

// dir with forward slashes
func AddDir(dir string) {
	fmt.Printf("D %v\n", dir)
	check(filepath.Walk(filepath.FromSlash(dir),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			AddFile(filepath.ToSlash(path))
			return nil
		}))
}
