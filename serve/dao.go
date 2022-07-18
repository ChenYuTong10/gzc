package main

import (
	mapset "github.com/deckarep/golang-set/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type DocDao struct {
	Coll  *mongo.Collection
	Index mongo.IndexModel
}

func NewDocDao() *DocDao {
	d := new(DocDao)
	d.Coll = mdb.Database(property.DB.DBName).Collection(property.DB.DocColl)
	d.Index = mongo.IndexModel{Keys: bson.D{{"Hash", 1}}}
	if _, err := d.Coll.Indexes().CreateOne(ctx, d.Index); err != nil {
		panic(err)
	}
	return d
}

func (d *DocDao) InsertMany(docs []Document) error {
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
	_, err := d.Coll.InsertMany(ctx, data)
	if err != nil {
		log.Println("unexpected error when insert document collection:", err, "docs:", docs, "data:", data)
		return err
	}
	return nil
}

func (d *DocDao) MatchHead(text string, num, size int64) ([]Document, error) {
	filter := bson.D{{"Head", bson.D{{"$regex", text}}}}
	option := options.Find().SetLimit(size).SetSkip((num - 1) * size).SetProjection(
		bson.D{{"Head", 1}, {"_id", 0}})

	cursor, err := d.Coll.Find(ctx, filter, option)
	if err != nil {
		log.Println("unexpected error when querying head with regex:", err, "text:", text)
		return nil, err
	}
	var docs []Document
	err = cursor.All(ctx, &docs)
	if err != nil {
		log.Println("unable to decode query document:", err, "docs:", docs, "text:", text)
	}
	return docs, nil
}

func (d *DocDao) Find(hash []string, num, size int64) ([]Document, error) {
	filter := bson.D{{"Hash", bson.D{{"$in", hash}}}}
	option := options.Find().SetLimit(size).SetSkip((num - 1) * size).SetProjection(
		bson.D{{"_id", 0}})

	cursor, err := d.Coll.Find(ctx, filter, option)
	if err != nil {
		log.Println("unexpected error when querying document with hash:", err, "hash:", hash)
		return nil, err
	}
	var docs []Document
	err = cursor.All(ctx, &docs)
	if err != nil {
		log.Println("unable to decode query document:", err, "docs:", docs, "hash:", hash)
	}
	return docs, nil
}

type AuthorDao struct {
	Coll  *mongo.Collection
	Index mongo.IndexModel
}

func NewAuthorDao() *AuthorDao {
	a := new(AuthorDao)
	a.Coll = mdb.Database(property.DB.DBName).Collection(property.DB.AuthorColl)
	a.Index = mongo.IndexModel{Keys: bson.D{{"Id", 1}}}
	if _, err := a.Coll.Indexes().CreateOne(ctx, a.Index); err != nil {
		panic(err)
	}
	return a
}

func (a *AuthorDao) InsertMany(authors []Author) error {
	var data []any
	for _, author := range authors {
		data = append(data, bson.D{
			{"Id", author.Id},
			{"Name", author.Name},
			{"Major", author.Major},
			{"Class", author.Class},
		})
	}
	_, err := a.Coll.InsertMany(ctx, data)
	if err != nil {
		log.Println("unexpected error when insert author collection:", err, "authors:", authors, "data:", data)
		return err
	}
	return nil
}

type WordDao struct {
	Coll  *mongo.Collection
	Index mongo.IndexModel
}

func NewWordDao() *WordDao {
	w := new(WordDao)
	w.Coll = mdb.Database(property.DB.DBName).Collection(property.DB.WordColl)
	w.Index = mongo.IndexModel{Keys: bson.D{{"Text", 1}}}
	if _, err := w.Coll.Indexes().CreateOne(ctx, w.Index); err != nil {
		panic(err)
	}
	return w
}

func (w *WordDao) InsertMany(words []Word) error {
	var data []any
	for _, word := range words {
		word.Json, _ = word.Docs.MarshalJSON()
		data = append(data, bson.D{
			{"Text", word.Text},
			{"Docs", word.Json},
		})
	}
	_, err := w.Coll.InsertMany(ctx, data)
	if err != nil {
		log.Println("unexpected error when insert word collection:", err, "words:", words, "data:", data)
		return err
	}
	return nil
}

func (w *WordDao) FindIndex(text string) (Word, error) {
	filter := bson.D{{"Text", text}}
	option := options.FindOne().SetProjection(
		bson.D{{"Docs", 1}, {"_id", 0}})

	word := Word{Text: text, Docs: mapset.NewSet[string]()}
	err := w.Coll.FindOne(ctx, filter, option).Decode(&word)
	if err != nil {
		log.Println("unable to decode query word:", err, "text:", text)
		return word, err
	}
	if err = word.Docs.UnmarshalJSON(word.Json); err != nil {
		return word, err
	}
	return word, nil
}
