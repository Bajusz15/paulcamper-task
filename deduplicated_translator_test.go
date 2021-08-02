package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"sync"
	"testing"
	"time"
)

type mockTranslatorForDeduplicated struct {
	mock.Mock
}

func (mt *mockTranslatorForDeduplicated) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	//return "mock translator result", nil
	args := mt.Called(ctx, from, to, data)
	time.Sleep(500 * time.Millisecond)
	return args.String(0), args.Error(1)
}

func Test_decoupledTranslator_Translate(t *testing.T) {
	m := new(mockTranslatorForDeduplicated)
	c := context.TODO()
	from := language.Hungarian
	to := language.English
	data := "example"
	dt := NewDeduplicatedTranslator(m)

	//testError := errors.New("this is a test error")
	var first time.Time
	var second time.Time
	m.On("Translate", c, from, to, data).Return("mock translator result", nil).Once().Run(func(args mock.Arguments) {
		first = time.Now()
	})
	m.On("Translate", c, from, to, data).Return("mock translator result", nil).Once().Run(func(args mock.Arguments) {
		second = time.Now()
	})

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(dt *deduplicatedTranslator, wg *sync.WaitGroup) {
		defer wg.Done()
		_, err := dt.Translate(c, from, to, data)
		assert.Equal(t, nil, err)
	}(dt, wg)

	//_, err := dt.Translate(c, from, to, data)
	//assert.Equal(t, testError, err)
	wg.Add(1)
	go func(dt *deduplicatedTranslator, wg *sync.WaitGroup) {
		defer wg.Done()
		_, err := dt.Translate(c, from, to, data)
		assert.Equal(t, nil, err)
	}(dt, wg)

	wg.Wait()

	diff := second.Sub(first)
	fmt.Println(diff)
	assert.Greater(t, diff, 500*time.Millisecond)
}
