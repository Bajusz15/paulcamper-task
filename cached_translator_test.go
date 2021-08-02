package main

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"testing"
)

type mockTranslatorForCache struct {
	mock.Mock
}

func (mt *mockTranslatorForCache) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	//return "mock translator result", nil
	args := mt.Called(ctx, from, to, data)
	return args.String(0), args.Error(1)
}

func Test_cachedTranslator_Translate(t *testing.T) {
	m := new(mockTranslatorForCache)
	c := context.TODO()
	from := language.Hungarian
	to := language.English
	data := "example"
	cache, _ := lru.New(100)
	cache.Add(from.String()+"-"+to.String()+"-"+data, "mock translator result")
	ct := &cachedTranslator{
		translator: m,
		repo:       cache,
	}
	_, err := ct.Translate(c, from, to, data)
	assert.Equal(t, nil, err)
	//testError := errors.New("this is a test error")
	m.AssertNotCalled(t, "Translate", c, from, to, data)
}
