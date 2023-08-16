package sensetive

import (
	"github.com/importcjj/sensitive"
	"log"
)

var Filter *sensitive.Filter

const WordDictPath = "./sensetive_doc/sensitiveWord.txt"

func InitFilter() {
	Filter = sensitive.New()
	err := Filter.LoadWordDict(WordDictPath)
	if err != nil {
		log.Println("InitFilter Fail,Err=" + err.Error())
	}
}
