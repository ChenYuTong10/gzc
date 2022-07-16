package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type DocRepo struct {
	Coll *mongo.Collection
}

func (dp *DocRepo) Init() *DocRepo {
	dp.Coll = mdb.Database(property.DB.DBName).Collection(property.DB.DocColl)
	return dp
}

func (dp *DocRepo) InsertMany(docs []Document) error {
	var data []any
	for _, doc := range docs {
		data = append(data, bson.D{
			{"Hash", doc.Hash},
			{"Size", doc.Size},
			{"AutId", doc.AutId},
			{"Order", doc.Order},
			{"Style", doc.Style},
			{"Source", doc.Source},
			{"Head", doc.Head},
			{"Body", doc.Body},
		})
	}
	_, err := dp.Coll.InsertMany(ctx, data)
	if err != nil {
		log.Println("unexpected error when insert document collection:", err, "docs:", docs, "data:", data)
		return err
	}
	return nil
}

type AuthorRepo struct {
	Coll *mongo.Collection
}

func (ar *AuthorRepo) Init() *AuthorRepo {
	ar.Coll = mdb.Database(property.DB.DBName).Collection(property.DB.AuthorColl)
	return ar
}

func (ar *AuthorRepo) InsertMany(authors []Author) error {
	var data []any
	for _, author := range authors {
		data = append(data, bson.D{
			{"_id", author.Id},
			{"Name", author.Name},
			{"Major", author.Major},
			{"Class", author.Class},
		})
	}
	_, err := ar.Coll.InsertMany(ctx, data)
	if err != nil {
		log.Println("unexpected error when insert author collection:", err, "authors:", authors, "data:", data)
		return err
	}
	return nil
}

type WordRepo struct {
	Coll *mongo.Collection
}

func (wr *WordRepo) Init() *WordRepo {
	wr.Coll = mdb.Database(property.DB.DBName).Collection(property.DB.WordColl)
	return wr
}

func (wr *WordRepo) InsertMany(words []Word) error {
	var data []any
	for _, word := range words {
		data = append(data, bson.D{
			{"Text", word.Text},
			{"Docs", word.Docs},
		})
	}
	_, err := wr.Coll.InsertMany(ctx, data)
	if err != nil {
		log.Println("unexpected error when insert word collection:", err, "words:", words, "data:", data)
		return err
	}
	return nil
}
