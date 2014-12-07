package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/russross/blackfriday"
)

// renderPage renders a given page by parsing Markdown for content
func renderPage(w http.ResponseWriter, r *http.Request, page string) {
	// Status code from flash cookies
	statusCode := http.StatusOK
	if page == "error" {
		statusCode = http.StatusInternalServerError
	}
	if fl := getSession(r).Flashes("error_code"); len(fl) > 0 {
		statusCode = fl[0].(int)
	}

	/* File reading */

	fileName := strings.ToLower(page) + ".md"
	if _, err := os.Stat(contentsPath + "/" + fileName); err != nil {
		redirectError(w, r, http.StatusNotFound)
		return
	}

	fileContents, err := ioutil.ReadFile(contentsPath + "/" + fileName)
	if err != nil {
		redirectError(w, r, http.StatusInternalServerError)
		return
	}

	// Prettifying the page name: "pAGe" --> "Page"
	displayedName := []byte(strings.ToLower(page))
	displayedName[0] = byte(unicode.ToUpper(rune(displayedName[0])))

	/* Markdown parsing and HTML rendering */

	// Error parsing
	pageContents := strings.Replace(
		string(fileContents),
		"{{.HTTP_STATUS_CODE}}",
		strconv.Itoa(statusCode),
		-1,
	)
	errorDescription := "That’s an error!"
	if fl := getSession(r).Flashes("error_description"); len(fl) > 0 {
		errorDescription = fl[0].(string)
	}
	pageContents = strings.Replace(pageContents, "{{.HTTP_ERROR_DESCRIPTION}}",
		errorDescription, -1)

	// Markdown parsing
	pageContents = string(blackfriday.Markdown(
		[]byte(pageContents),
		blackfriday.HtmlRenderer(0, "", ""),
		blackfriday.EXTENSION_LAX_HTML_BLOCKS,
	))

	// OCR result display (after the Markdown rendering, to add some HTML)
	ocrResult := ""
	if fl := getSession(r).Flashes("ocr_result"); len(fl) > 0 {
		ocrResult = fl[0].(string)
		ocrResultSplit := strings.Split(ocrResult, "\n")
		ocrResult = strings.Join(ocrResultSplit[1:len(ocrResultSplit)], "<br />")
		ocrResult = "Text recognised by our OCR:<br /><br />" + ocrResult
	}
	pageContents = strings.Replace(pageContents, "{{.OCR_RESULT}}", ocrResult,
		-1)

	saveSession(w, r)

	getRender(w, r).HTML(w, statusCode, "index", map[string]interface{}{
		"host":     r.Host,
		"title":    string(displayedName),
		"contents": template.HTML(pageContents),
	})
}
