package utils

import "github.com/yanyiwu/gojieba"

type Tokenizer struct {
	t *gojieba.Jieba
}

func (tz *Tokenizer) Init() *Tokenizer {
	tz.t = gojieba.NewJieba()
	return tz
}

func (tz *Tokenizer) CutForSearch(s string) []string {
	return tz.t.CutForSearch(s, true)
}
