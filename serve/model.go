package main

import (
	u "serve/utils"
)

type Style int

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

//+-------+----------+
//| Genre |   Type   |
//+-------+----------+
//|   S   |   散文   |
//|   K   |   小说   |
//|   P   |   诗歌   |
//|   J   |   剧本   |
//|   Y   |   评论   |
//|   N   |   新闻   |
//|   R   |   论文   |
//|   O   | 自我介绍 |
//|   I   |  应用文  |
//+-------+----------+

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
	Hash   string `json:"hash"` // body hash value
	Size   int64  `json:"size"`
	AutId  string `json:"authorId"` // author id
	Order  int64  `json:"order"`    // written sequence
	Style  Style  `json:"style"`    // type of literature
	Source string `json:"source"`
	Head   string `json:"head"`
	Body   string `json:"body"`
}

type Author struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Major string `json:"major"`
	Class string `json:"class"`
}

type Word struct {
	Text string        `json:"text"`
	Docs u.Set[string] `json:"docs"`
}
