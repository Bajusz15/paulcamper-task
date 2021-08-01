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

func (ct *cachedTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	val, ok := ct.repo.Get(from.String() + "-" + to.String() + "-" + data)
	if ok {
		return val.(string), nil
	}

	result, err := ct.translator.Translate(ctx, from, to, data)
	if err != nil {
		return "", err
	}

	ct.repo.Add(from.String()+"-"+to.String()+"-"+data, result)

	return result, nil
}
