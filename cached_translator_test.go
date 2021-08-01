package main

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/text/language"
	"testing"
)

func Test_cachedTranslator_Translate(t *testing.T) {
	type fields struct {
		translator Translator
		repo       *lru.Cache
	}
	type args struct {
		ctx  context.Context
		from language.Tag
		to   language.Tag
		data string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &cachedTranslator{
				translator: tt.fields.translator,
				repo:       tt.fields.repo,
			}
			got, err := ct.Translate(tt.args.ctx, tt.args.from, tt.args.to, tt.args.data)
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
