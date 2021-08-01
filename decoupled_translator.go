package main

import (
	"context"
	"golang.org/x/text/language"
	"runtime"
	"sync"
)

type decoupledTranslator struct {
	translator Translator
	requestMap map[string]bool
	mux        *sync.RWMutex
}

func (dt *decoupledTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	key := from.String() + "-" + to.String() + "-" + data
	dt.mux.RLock()
	for dt.requestMap[key] == true {
		//TODO: wait here until it is false
		dt.mux.RUnlock()
		runtime.Gosched()
		dt.mux.RLock()
	}
	dt.mux.RUnlock()

	dt.mux.Lock()
	dt.requestMap[key] = true
	dt.mux.Unlock()
	result, err := dt.translator.Translate(ctx, from, to, data)
	dt.mux.Lock()
	dt.requestMap[key] = false
	dt.mux.Unlock()
	if err != nil {
		return "", err
	}

	return result, nil
}
