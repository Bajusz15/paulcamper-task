package main

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func Test_cachedTranslator_Translate(t *testing.T) {
	m := new(mockTranslator)
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
