package main

import (
	lru "github.com/hashicorp/golang-lru"
	"log"
	"time"
)

// Service is a Translator user.
type Service struct {
	translator Translator
}

func NewService() *Service {

	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)
	cache, err := lru.New(100)
	if err != nil {
		log.Fatalln("could not create cache")
	}
	bt := NewBackoffTranslator(t, 16*time.Second, 10)
	ct := NewCachedTranslator(bt, cache)
	dt := NewDeduplicatedTranslator(ct)

	return &Service{
		translator: dt,
	}
}

//TODO: task 1 exponential backoff
//TODO: task 2 memory repository. save "from", "to", "test" parameters so we know if same query is same. Then we just use cache instead of making the request
//TODO: task 3 same thing as before, but use mutex instead of a repository to not make the same request again
