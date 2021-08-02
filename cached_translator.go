package main

import (
	"context"
	"github.com/hashicorp/golang-lru"
	"golang.org/x/text/language"
)

type cachedTranslator struct {
	translator Translator
	repo       *lru.Cache
}

func NewCachedTranslator(t Translator, cache *lru.Cache) *cachedTranslator {
	return &cachedTranslator{t, cache}
}

func (ct *cachedTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	key := from.String() + "-" + to.String() + "-" + data
	val, ok := ct.repo.Get(key)
	if ok {
		return val.(string), nil
	}

	result, err := ct.translator.Translate(ctx, from, to, data)
	if err != nil {
		return "", err
	}

	ct.repo.Add(key, result)

	return result, nil
}
