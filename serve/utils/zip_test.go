package utils

import (
	"log"
	"testing"
)

func TestDecompress(t *testing.T) {
	if err := Decompress("D:\\Code\\2022挑战杯语料库\\serve\\raw\\2018级千字文 原文.zip"); err != nil {
		log.Println("decompress error:", err)
		return
	}
	log.Println("decompress ok")
}
