package utils

import (
	"log"
	"testing"
)

func TestNewSet(t *testing.T) {
	s1 := NewSet[int]()

	// add
	s1.Add(1)
	log.Println("s1 add 1:", s1)
	s1.Add(2)
	log.Println("s1 add 2:", s1)

	// range
	log.Println("s1 range:", s1.Range())

	// delete
	s1.Del(1)
	log.Println("s1 delete 1:", s1)
	s1.Del(3)
	log.Println("s1 delete 3:", s1)

	// find
	ok := s1.Find(2)
	log.Println("s1 find 2:", ok)
	ok = s1.Find(3)
	log.Println("s1 find 3:", ok)
}
