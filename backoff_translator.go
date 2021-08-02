package main

import (
	"context"
	"fmt"
	"github.com/pailcamper/pc-offline-challenge/pkg/backoff"
	"golang.org/x/text/language"
	"time"
)

type backoffTranslator struct {
	translator Translator

	backoffService *backoff.Service
}

func NewBackoffTranslator(t Translator, maxBackoff time.Duration, retries int) *backoffTranslator {
	return &backoffTranslator{t, backoff.NewService(maxBackoff, retries)}
}

func (bt *backoffTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	var result string
	var err error

	bt.backoffService.Try(func() error {
		result, err = bt.translator.Translate(ctx, from, to, data)
		fmt.Println(err)
		return err
	})

	if err != nil {
		return "", err
	}

	return result, nil
}
