package main

import (
	lru "github.com/hashicorp/golang-lru"
	"log"
	"sync"
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
	var decoupleMap map[string]bool
	var mutex sync.RWMutex
	return &Service{
		translator: &deduplicatedTranslator{
			translator: &cachedTranslator{
				translator: newBackoffTranslator(t, 10*time.Second, 5),
				repo:       cache,
			},
			requestMap: decoupleMap,
			mux:        &mutex,
		},
	}
}

//TODO: task 1 exponential backoff
//TODO: task 2 memory repository. save "from", "to", "test" parameters so we know if same query is same. Then we just use cache instead of making the request
//TODO: task 3 same thing as before, but use mutex instead of a repository to not make the same request again
