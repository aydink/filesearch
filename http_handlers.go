package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/gabriel-vasile/mimetype"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index", nil)
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {

	file := r.URL.Query().Get("file")
	//fmt.Println("Dosya:", file)

	pathSlice := strings.Split(file, string(os.PathSeparator))
	fileName := pathSlice[len(pathSlice)-1]

	infile, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(w, "Dosya bunamadı: %s", file)
		return
	}

	mime, err := mimetype.DetectReader(infile)

	//fileHeader := make([]byte, 512)

	//Copy the headers into the fileHeader buffer

	//infile.Read(fileHeader)
	//Get content type of file
	//fileContentType := http.DetectContentType(fileHeader)

	//Get the file size
	fileStat, _ := infile.Stat()                       //Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", mime.String())
	w.Header().Set("Content-Length", fileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	infile.Seek(0, 0)
	io.Copy(w, infile) //'Copy' the file to the client

	if err != nil {
		fmt.Println("cannot open file")
	}

	io.Copy(w, infile)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{}

	order := r.URL.Query().Get("order")
	q := r.URL.Query().Get("q")
	data.Query = q

	q = strings.TrimSpace(q)
	q = strings.ToLowerSpecial(unicode.TurkishCase, q)
	queryTerms := strings.Split(q, " ")

	var searchResult []FileMeta

	for _, fileMeta := range destinationMeta {
		hasAllKeys := true

		for _, term := range queryTerms {
			if term != "" {
				if !strings.Contains(strings.ToLowerSpecial(unicode.TurkishCase, fileMeta.FileName), term) {
					hasAllKeys = false
				}
			}
		}

		if hasAllKeys {
			searchResult = append(searchResult, fileMeta)
		}
	}

	if order == "mtime" {
		sort.Sort(ByModTime(searchResult))
	}

	if order == "name" {
		sort.Sort(ByName(searchResult))
	}

	if order == "size" {
		sort.Sort(BySize(searchResult))
	}

	data.NumFiles = len(searchResult)
	if len(searchResult) > 2000 {
		data.Truncated = true
		data.Result = searchResult[:2000]
	} else {
		data.Truncated = false
		data.Result = searchResult
	}

	err := tpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		fmt.Println(err)
	}
}
