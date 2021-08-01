package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"testing"
	"time"
)

type mockTranslator struct {
	mock.Mock
}

func (mt *mockTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	//return "mock translator result", nil
	args := mt.Called(ctx, from, to, data)
	return args.String(0), args.Error(1)
}

func Test_backoffTranslator_Retries(t *testing.T) {
	m := new(mockTranslator)
	c := context.TODO()
	from := language.Hungarian
	to := language.English
	data := "example"
	bt := newBackoffTranslator(m, 2*time.Second, 10)

	start := time.Now()
	testError := errors.New("this is a test error")
	var second time.Time
	m.On("Translate", c, from, to, data).Return("", testError).Twice()
	_, err := bt.Translate(c, from, to, data)
	assert.Equal(t, testError, err)
	_, err = bt.Translate(c, from, to, data)
	assert.Equal(t, testError, err)
	m.On("Translate", c, from, to, data).Return("mock translator result", nil).Once().Run(func(args mock.Arguments) {
		second = time.Now()
	})

	//_, err := bt.Translate(c, from, to, data)
	//assert.Equal(t, testError, err)
	_, err = bt.Translate(c, from, to, data)
	assert.Equal(t, nil, err)
	diff := second.Sub(start)
	fmt.Println(diff)
	assert.Greater(t, diff, 500*time.Millisecond)
}

func Test_backoffTranslator_Error(t *testing.T) {
	m := new(mockTranslator)
	c := context.TODO()
	from := language.Hungarian
	to := language.English
	data := "example"
	testError := errors.New("this is a test error")

	m.On("Translate", c, from, to, data).Return("", testError)
	bt := newBackoffTranslator(m, 2*time.Second, 10)

	result, err := bt.Translate(c, from, to, data)
	assert.Equal(t, testError, err)
	assert.Equal(t, "", result)
}

func Test_backoffTranslator_Translate(t *testing.T) {
	type fields struct {
		translator Translator
	}
	type args struct {
		ctx  context.Context
		from language.Tag
		to   language.Tag
		data string
	}
	m := new(mockTranslator)
	bt := newBackoffTranslator(m, 2*time.Second, 10)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "test no error",
			fields: fields{translator: bt},
			args: args{
				ctx:  context.TODO(),
				from: language.Hungarian,
				to:   language.English,
				data: "example",
			},
			want:    "mock translator result",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		m.On("Translate", tt.args.ctx, tt.args.from, tt.args.to, tt.args.data).Return("mock translator result", nil)
		t.Run(tt.name, func(t *testing.T) {
			bf := backoffTranslator{
				translator: tt.fields.translator,
			}
			got, err := bf.Translate(tt.args.ctx, tt.args.from, tt.args.to, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Translate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
