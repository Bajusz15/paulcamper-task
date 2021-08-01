package main

import (
	"context"
	"github.com/pailcamper/pc-offline-challenge/pkg/backoff"
	"golang.org/x/text/language"
)

type backoffTranslator struct {
	translator Translator
}

func (bf backoffTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	backoffService := backoff.NewService()
	var result string
	var err error

	backoffService.Try(func() error {
		result, err = bf.translator.Translate(ctx, from, to, data)
		return err
	})

	return result, nil
}
