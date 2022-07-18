package main

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type Style = int

const (
	News Style = iota
	Play
	Poem
	Novel
	Paper
	Prose
	Review
	Practical
	SelfIntro
)

// StyleTb stores the relation of between symbol and its style.
var StyleTb = map[string]Style{
	"N": News,
	"J": Play,
	"P": Poem,
	"K": Novel,
	"R": Paper,
	"S": Prose,
	"Y": Review,
	"I": Practical,
	"O": SelfIntro,
}

type Document struct {
	Hash   string // body hash value
	Size   int64
	AutId  string // author id
	Order  int64  // written sequence
	Style  Style  // type of literature
	Source string
	Head   string
	Body   string
}

type Author struct {
	Id    string
	Name  string
	Major string
	Class string
}

type Word struct {
	Text string
	Json []byte `bson:"Docs"` // the string of docs
	Docs mapset.Set[string]
}
