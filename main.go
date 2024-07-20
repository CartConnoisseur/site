package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var pageFormat string
var errorPage []byte

func render(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if len(path) == 0 {
		path = "index"
	}

	page, err := os.ReadFile("pages/" + path + ".md")
	if err != nil {
		page = errorPage
	}

	title := path

	if page[0] == '#' {
		title = strings.Split(string(page), "\n")[0][2:]
	}

	fmt.Fprintf(w, pageFormat, title, string(render(page)))
}

func main() {
	bytes, err := os.ReadFile("page.html")
	if err != nil {
		panic(err)
	}
	pageFormat = string(bytes)

	errorPage, err = os.ReadFile("missing.md")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, r.URL.Path[1:]) })
	log.Fatal(http.ListenAndServe(":8080", nil))
}
