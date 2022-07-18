package main

import (
	"context"
	"fmt"
	"github.com/ChenYuTong10/ini"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type Property struct {
	Path struct {
		Pro string `ini:"pro"`
		Raw string `ini:"raw"`
		Web string `ini:"web"`
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
	mdb       *mongo.Client
	ctx       context.Context
	property  *Property
	docDao    *DocDao
	wordDao   *WordDao
	authorDao *AuthorDao
)

func init() {
	property = new(Property)
	if err := ini.Bind("application.ini", property); err != nil {
		log.Fatalln(err)
	}

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
	docDao = NewDocDao()
	wordDao = NewWordDao()
	authorDao = NewAuthorDao()
}

func main() {
	//BatchDir(filepath.Join(property.Path.Pro, property.Path.Raw, "2018级千字文 原文"))
	http.HandleFunc("/", Gateway)
	if err := http.ListenAndServe(property.Serve.Port, nil); err != nil {
		log.Fatalln("unable to start service:", err)
	}
}
