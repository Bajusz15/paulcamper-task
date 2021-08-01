package main

import (
	"context"
	"golang.org/x/text/language"
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
	if dt.requestMap[key] {
		//TODO: wait here until it is false
	}
	dt.mux.RUnlock()

	dt.mux.Lock()
	result, err := dt.translator.Translate(ctx, from, to, data)
	dt.mux.Unlock()
	if err != nil {
		return "", err
	}

	return result, nil
}
