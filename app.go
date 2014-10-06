package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/russross/blackfriday"
)

const (
	ContentsPath = "./contents"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		renderPage(200, "home", r)
	})

	m.Get("/:page", func(params martini.Params, r render.Render) {
		renderPage(200, params["page"], r)
	})

	m.NotFound(func(r render.Render) {
		renderPage(404, "error", r)
	})

	m.Run()
}

// renderPage renders a given page (without trusting user input)
func renderPage(statusCode int, pageName string, r render.Render) {
	if checkPageName(pageName) != nil {
		renderPage(404, "error", r)
		return
	}
	fileName := strings.ToLower(pageName) + ".md"
	if _, err := os.Stat(ContentsPath + "/" + fileName); err != nil {
		renderPage(404, "error", r)
		return
	}

	displayedName := []byte(strings.ToLower(pageName))
	displayedName[0] = byte(unicode.ToUpper(rune(displayedName[0])))

	fileContents, err := ioutil.ReadFile(ContentsPath + "/" + fileName)
	if err != nil {
		renderPage(500, "error", r)
		return
	}

	pageContentsRaw := string(blackfriday.MarkdownCommon(fileContents))
	pageContentsRaw = strings.Replace(pageContentsRaw, "{{.HTTP_STATUS_CODE}}",
		strconv.Itoa(statusCode), -1)
	pageContents := template.HTML(pageContentsRaw)

	r.HTML(statusCode, "index", map[string]interface{}{
		"title":    string(displayedName),
		"contents": pageContents,
	})
}

// checkPageName returns an error if the name of a page ia invalid (ie. insecure)
func checkPageName(pageName string) error {
	for _, c := range pageName {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return errors.New("Invalid character found inside page name")
		}
	}
	return nil
}