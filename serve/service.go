package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	u "serve/utils"
	"strings"
)

var uriRegx = regexp.MustCompile("^/(?P<action>\\w+)/(?P<pagenum>\\w+)/(?P<keyword>\\w+)(\\?(?P<query>.*))?")

const (
	action = iota
	pagenum
	keyword
)

const UriMatchLength = 3

const (
	Method  = "method"
	Keyword = "keyword"
)

func Gateway(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.URL.Query(), r.Body)

	//secs := uriRegx.FindStringSubmatch(r.URL.RequestURI())
	//if len(secs) < UriMatchLength {
	//	http.Error(w, "403 Forbidden", http.StatusForbidden)
	//	return
	//}
	switch r.URL.Path {
	case "/search":
		Search(w, r)
	default:
		http.ServeFile(w, r, "")
		//http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

func FileServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.FileServer(http.Dir(filepath.Join(property.Path.Pro, property.Path.Web))).ServeHTTP(w, r)
	})
}

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	remote, header, err := r.FormFile("chunk")
	if err != nil {
		log.Println("unable to accept request file:", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	if filepath.Ext(header.Filename) != ZipSuffix {
		if _, err = w.Write([]byte("unsupported type")); err != nil {
			log.Println("unexpected error when writing back:", err)
		}
	}
	zipPath := filepath.Join(property.Path.Pro, property.Path.Raw, header.Filename)

	local, err := os.Create(zipPath)
	if err != nil {
		log.Println("unable to create local file:", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(local, remote)
	if err != nil {
		log.Println("unable to save request file:", err, "save path:", zipPath)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Why don't we use defer?
	// The following process, such as Decompress and BatchDir, needs to handle the file.
	// If the file is not closed, other process can not get the permission to handle it.
	if err = local.Close(); err != nil {
		log.Println("unexpected error when closing file:", err, "save path:", zipPath)
	}
	go func() {
		if err = u.Decompress(zipPath); err != nil {
			log.Println("unable to decompress request zip:", err, "absolute path:", zipPath)
			return
		}
		// The name of decompressed directory is same as the zip file expected extension.
		BatchDir(filepath.Join(property.Path.Pro, property.Path.Raw, strings.Split(header.Filename, Dot)[0]))
	}()

	if _, err = w.Write([]byte("ok")); err != nil {
		log.Println("unexpected error when writing back:", err)
	}
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()

	text := q.Get("text")
	group := q.Get("group")
	limit := q.Get("limit")
	genre := q.Get("genre")
	grade := q.Get("grade")

	if u.HasEmpty(text, group, limit, genre, grade) {
		// return no data
	}

	var err error
	var docs []Document
	if group == "Head" {
		docs, err = docDao.MatchHead(text, 1, 10)
		if err != nil {
			log.Println("query document error:", err, "q:", q)
		}
	} else {
		var word Word
		word, err = wordDao.FindIndex(text)
		if err != nil {
			log.Println("query word docs error:", err, "text", text, "q:", q)
		}
		docIndexes := u.ConvertSetToArray[string](word.Docs)
		docs, err = docDao.Find(docIndexes, 1, 10)
		if err != nil {
			log.Println("query docs error:", err, "text", text, "word:", word, "q:", q)
		}
	}
	log.Println(docs, len(docs))
}
