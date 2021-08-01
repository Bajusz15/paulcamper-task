package main

import (
	"context"
	"github.com/pailcamper/pc-offline-challenge/pkg/backoff"
	"golang.org/x/text/language"
	"time"
)

type backoffTranslator struct {
	translator Translator
	maxBackoff time.Duration
}

func (bt *backoffTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	backoffService := backoff.NewService(bt.maxBackoff)
	var result string
	var err error

	backoffService.Try(func() error {
		result, err = bt.translator.Translate(ctx, from, to, data)
		return err
	})

	return result, nil
}
