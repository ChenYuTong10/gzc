package main

import (
	"log"
	"testing"
)

func TestRegx(t *testing.T) {
	p1 := "2018级千字文 原文\\汉语182班 原文\\B1863110023林浩\\B1863110023林浩.TXT\\B1863110023S16.txt"
	r1 := metaRegx.FindStringSubmatch(p1)
	log.Println("p1 results:", r1, "length:", len(r1))
	p2 := "2018级千字文 原文\\汉语181班 原文\\A1812100048 梁坚玲\\A1812100048 梁坚玲 TXT\\A1812100048K01.txt"
	r2 := metaRegx.FindStringSubmatch(p2)
	log.Println("p2 results:", r2, "length:", len(r2))
	p3 := "D:\\Code\\2022挑战杯语料库\\serve\\raw\\2018级千字文 原文\\汉语182班 原文\\B1812100076曾宏晟\\B1812100076曾宏晟TXT\\B1812100076O .txt"
	r3 := metaRegx.FindStringSubmatch(p3)
	log.Println("p3 results:", r3, "length:", len(r3))
}
