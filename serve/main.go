package main

import (
	"context"
	"fmt"
	"github.com/ChenYuTong10/ini"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"path/filepath"
)

type Property struct {
	Path struct {
		Pro string `ini:"pro"`
		Raw string `ini:"raw"`
	} `ini:"path"`
	Serve struct {
		Port string `ini:"port"`
	} `ini:"serve"`
	DB struct {
		Addr       string `ini:"addr"`
		DBName     string `ini:"dbname"`
		DocColl    string `ini:"docColl"`
		WordColl   string `ini:"wordColl"`
		AuthorColl string `ini:"authorColl"`
		SourceColl string `ini:"sourceColl"`
		Passwd     string `ini:"passwd"`
		User       string `ini:"user"`
	} `ini:"db"`
}

var (
	property *Property

	mdb        *mongo.Client
	ctx        context.Context
	docRepo    *DocRepo
	wordRepo   *WordRepo
	authorRepo *AuthorRepo
)

func init() {
	property = new(Property)
	if err := ini.Bind("application.ini", property); err != nil {
		log.Fatalln(err)
	}

	// repository
	ctx = context.Background()
	var err error
	mdb, err = mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf(
			"mongodb://%v:%v@%v/?maxPoolSize=20&w=majority",
			property.DB.User,
			property.DB.Passwd,
			property.DB.Addr,
		)),
	)
	if err != nil {
		log.Fatalln("unable to establish connection with mongodb:", err)
	}
	if err = mdb.Ping(ctx, nil); err != nil {
		log.Fatalln("unsure with mongodb connection alive :", err)
	}
	docRepo = new(DocRepo).Init()
	authorRepo = new(AuthorRepo).Init()
	wordRepo = new(WordRepo).Init()
}

func main() {
	//log.SetFlags(log.Llongfile)
	//http.HandleFunc("/", Upload)
	//http.ListenAndServe(property.Serve.Port, nil)
	BatchDir(filepath.Join(property.Path.Pro, property.Path.Raw, "2018级千字文 原文"))
}
