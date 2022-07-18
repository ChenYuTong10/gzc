package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/ChenYuTong10/chardet"
	mapset "github.com/deckarep/golang-set/v2"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	u "serve/utils"
	"strconv"
)

const (
	HeadPrefix = "BBB"
	HeadSuffix = "TTT"

	// MetaMatchLength includes origin, source, major, class, id, name, style and order.
	MetaMatchLength = 8
)

const (
	origin = iota
	source
	major
	class
	autId // author id
	autNa // author name
	style
	order
)

var (
	headRegx  = regexp.MustCompile("第.*篇|[《》]*")
	bodyRegx  = regexp.MustCompile("[\a\f\t\n\r\v ]+")
	metaRegx  = regexp.MustCompile("(?P<source>.*?)[\\\\/](?P<major>[\u4e00-\u9fa5]+)[\t ]*(?P<seq>[0-9]{3}).*?[a-zA-Z](?P<id>[0-9]{10})[\t ]?(?P<name>[\u4e00-\u9fa5]+).*?(?P<style>[a-zA-Z])(?P<order>[0-9]+)?[\t ]?.txt")
	tokenizer = new(u.Tokenizer).Init()
)

var (
	ErrPathFormat = errors.New("specific path is not formatted")
	ErrNoDocHead  = errors.New("document head is not exist")
	ErrNoDocBody  = errors.New("document body is not exist")
)

func BatchDir(dir string) {
	var lastDir string
	var docs []Document
	var authors []Author
	var wordMap = make(map[string]mapset.Set[string])
	var badPath []string // the path of bad file(e.g: not formatted, encode wrongly)
	_ = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Println("unable to visit file:", path, err)
			// filepath.Walk will be stopped once returning a non-skip error.
			// So returning a nil error to maintain the walk function.
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == TextSuffix {
			doc, author, cuts, err := BatchFile(path, info)
			if err != nil {
				badPath = append(badPath, path)
				log.Println("batch file error:", err, "dir:", dir, "path:", path)
				return nil
			}
			docs = append(docs, *doc)
			if d := filepath.Dir(path); d != lastDir {
				lastDir = d
				authors = append(authors, *author)
			}
			for _, text := range cuts {
				set, loaded := wordMap[text]
				if !loaded {
					set = mapset.NewSet[string]()
				}
				set.Add(doc.Hash)
				wordMap[text] = set
			}
		}
		return nil
	})
	if err := docDao.InsertMany(docs); err != nil {
		log.Println("storing documents error:", err, "docs:", docs)
		return
	}
	if err := authorDao.InsertMany(authors); err != nil {
		log.Println("storing authors error:", err, "authors:", authors)
		return
	}
	var words []Word
	for text, set := range wordMap {
		words = append(words, Word{Text: text, Docs: set})
	}
	if err := wordDao.InsertMany(words); err != nil {
		log.Println("storing words error:", err, "words:", words)
		return
	}
	log.Println("batch directory done", "dir:", dir)
}

func BatchFile(path string, info fs.FileInfo) (*Document, *Author, []string, error) {
	secs := metaRegx.FindStringSubmatch(path)
	if len(secs) < MetaMatchLength {
		// the directory or text file is not formatted.
		//unFormatted = append(unFormatted, path)
		log.Println("unable to match with path:", path, "match result:", secs)
		return nil, nil, nil, ErrPathFormat
	}
	author := &Author{
		Id:    secs[autId],
		Name:  secs[autNa],
		Major: secs[major],
		Class: secs[class],
	}

	stream, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("unexpected error when reading file:", err, "path:", path)
		return nil, author, nil, err
	}
	d := new(chardet.Detector)
	d.Feed(stream)
	if d.State && d.Encoding != chardet.UTF8 {
		stream, err = u.EncodeUTF8(stream, d.Encoding)
		if err != nil {
			log.Println("unexpected error when encoding stream:", err, "path:", path)
			return nil, author, nil, err
		}
		if err = ioutil.WriteFile(path, stream, 0666); err != nil {
			log.Println("unable to write stream back to file:", err, "path:", path)
			return nil, author, nil, err
		}
	}

	_, head, body := u.Cut(string(stream), HeadPrefix, HeadSuffix)
	if len(head) == 0 {
		log.Println("document head can not be empty", "path:", path)
		return nil, author, nil, ErrNoDocHead
	}
	if len(body) == 0 {
		log.Println("document body can not be empty", "path:", path)
		return nil, author, nil, ErrNoDocBody
	}
	head = headRegx.ReplaceAllString(head, "")
	body = bodyRegx.ReplaceAllString(body, "")
	hash := fmt.Sprintf("%x", md5.Sum([]byte(body)))

	order, err := strconv.ParseInt(secs[order], 10, 64)
	if err != nil {
		log.Println("invalid document order:", err, "path:", path, "match:", secs)
		return nil, author, nil, err
	}
	doc := &Document{
		Hash:   hash,
		Size:   info.Size(),
		AutId:  author.Id,
		Order:  order,
		Style:  StyleTb[secs[style]],
		Source: filepath.Base(secs[source]),
		Head:   head,
		Body:   body,
	}
	bodyWords := tokenizer.CutForSearch(doc.Body)

	return doc, author, bodyWords, nil
}
