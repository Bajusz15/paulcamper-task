package main

import (
	"context"
	"golang.org/x/text/language"
	"sync"
)

type deduplicatedTranslator struct {
	translator           Translator
	requestMap           map[string]bool
	mux                  *sync.Mutex
	resourceSynchronizer *sync.Cond
}

func NewDeduplicatedTranslator(t Translator) *deduplicatedTranslator {
	mutex := sync.Mutex{}
	condition := sync.NewCond(&mutex)
	return &deduplicatedTranslator{t, map[string]bool{}, &mutex, condition}
}

func (dt *deduplicatedTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	key := from.String() + "-" + to.String() + "-" + data
	dt.resourceSynchronizer.L.Lock()
	for dt.requestMap[key] == true {
		dt.resourceSynchronizer.Wait()
		//runtime.Gosched()
	}
	dt.resourceSynchronizer.L.Unlock()

	dt.resourceSynchronizer.L.Lock()
	dt.requestMap[key] = true
	dt.resourceSynchronizer.Broadcast()
	dt.resourceSynchronizer.L.Unlock()
	result, err := dt.translator.Translate(ctx, from, to, data)
	//dt.mux.Lock()
	dt.resourceSynchronizer.L.Lock()
	dt.requestMap[key] = false
	dt.resourceSynchronizer.Broadcast()
	dt.resourceSynchronizer.L.Unlock()
	if err != nil {
		return "", err
	}

	return result, nil
}
