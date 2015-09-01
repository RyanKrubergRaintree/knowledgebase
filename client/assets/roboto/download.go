package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const fonturl = `http://fonts.googleapis.com/css?family=RobotoDraft:300,500,700,400`

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func download(url string) []byte {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	check(err)

	return data
}

func main() {
	rx := regexp.MustCompile(`url\(([^)]+)\)`)
	index := 0
	origcss := download(fonturl)
	css := rx.ReplaceAllStringFunc(string(origcss),
		func(url string) string {
			url = strings.TrimPrefix(url, "url(")
			url = strings.TrimSuffix(url, ")")

			data := download(url)
			filename := fmt.Sprintf("font.%v.woff2", index)
			ioutil.WriteFile(filename, data, 0777)
			index++

			return "url(" + filename + ")"
		})
	ioutil.WriteFile("all.css", []byte(css), 0777)
}
